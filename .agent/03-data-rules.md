# 分册 7｜数据库与消息中间件规范

> 归属总纲：`../1-总纲/AGENT.md`
> 覆盖：PostgreSQL 建模、pgvector 向量检索、Kafka 消息、DDL 变更流程

---

## 0. 双库隔离（P0，优先于本分册其他规范）

| 库 | 句柄 | 存放内容 | 禁止 |
|----|------|----------|------|
| gva 库 | `global.GVA_DB` | gin-vue-admin 原生系统表（`sys_*`、文件、jwt、通知、字典、权限等框架内置表） | **禁止新建任何业务表、业务 model**；业务代码不得往 gva 库写业务数据 |
| luobo 库 | `global.LUOBO_DB` | 全部业务表（`cs_*` 客服、乘客、订单、钱包、优惠券等，含后续新增业务表） | 禁止把业务表建在 gva 库；禁止业务 model 使用 `GVA_DB` |

**表前缀归属**：`sys_*` → gva 库（只读，禁建业务表）；`ride_*` / `cs_*` / `dict_*` 及后续新增业务表 → luobo 库。

1. **严格隔离**：`GVA_DB` 只负责框架自带系统功能；所有业务 CRUD、业务表、业务 model 一律使用 `LUOBO_DB`；业务表禁止创建在 gva 库。
2. **`GVA_MODEL` 只是字段结构体（重要）**：业务 model **可以**嵌入 `global.GVA_MODEL`，但**仅为复用 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt` 四个字段**；`GVA_MODEL` **不会自动绑定任何数据库连接**。业务增删改查必须**显式使用 `global.LUOBO_DB` 操作 model 实例**，**禁止用 model 实例变量直接调用 GORM 方法**（Find / Create / Save / Updates / Delete 等），model 实例只作数据载体。
3. **事务统一在 `service` 层开启**：业务事务使用 `global.LUOBO_DB.Transaction`；**禁止多层事务嵌套**；事务内部禁止发起外部 HTTP 调用、Kafka 同步发送。
4. **只读统计查询**：允许基于 `global.LUOBO_DB` 开 GORM 只读 `Session`；**禁止**代码内单独 `gorm.Open` 新建 DB 连接、禁止手写 DSN。
5. **向量查询参数化**：pgvector 查询必须参数化，禁止字符串拼接向量数组（见 §2.4 第 4 条）。
6. **SQL 日志**：开发环境可开 GORM `Debug`；测试、生产环境**禁止打印 SQL 绑定参数**。
7. **链式起点必须是句柄**：所有 GORM 链式调用的**第一个接收者**必须是 `global.LUOBO_DB`（业务）或 `global.GVA_DB`（框架 `sys_*` 系统表）；**禁止以 model 结构体实例作为链式起点**。
8. **禁止在 model 上写 CRUD 方法**：不得为业务 model 结构体定义 Find / Create / Update / Delete / Save 等成员方法；数据库操作代码统一放在 `service` 层，`model` 只保留结构体、`TableName()`、字段 tag。
9. **软删除统一走 `Delete`**：用 `global.LUOBO_DB.Delete(&entity)`，由 GORM 自动维护 `DeletedAt`；**禁止手动 `Update` 修改 `deleted_at` 字段**；查询已删除数据必须显式调用 `Unscoped()`。

---

## 1. PostgreSQL 建模规范

### 1.1 通用约定

| 项 | 规范 |
|----|------|
| 主键 | 雪花 ID（`global.Snow.NextID()`），类型 `BIGINT`，**禁止自增主键** |
| 金额 | `decimal`，**禁止 float/double** |
| 时间 | `TIMESTAMPTZ(6)`；可空用指针字段 |
| 软删除 | `deleted_at TIMESTAMPTZ(6)`（由 `global.GVA_MODEL` 提供），**禁止自建删除标记字段** |
| 审计 | `created_by` / `updated_by BIGINT NOT NULL DEFAULT 0` |
| 注释 | 表与字段必须 `COMMENT`，与 Go model 的 `comment` tag 保持一致 |
| 索引 | 查询字段、外键、状态字段建索引；软删除表统一 `idx_xxx_deleted_at` |

> 软删除口径：统一 `global.LUOBO_DB.Delete(&entity)`，由 GORM 自动维护 `deleted_at`；**禁止手动 `Update` 修改 `deleted_at`**；默认查询自动排除已删除数据，**查已删除数据必须显式 `Unscoped()`**。

### 1.2 表前缀

| 前缀 | 含义 | 权限 |
|------|------|------|
| `sys_*` | GVA 系统自带表（用户/角色/菜单/字典…） | 业务只读，**禁止新增 `sys_` 业务表** |
| `ride_*` | 打车业务核心表（订单、司机、乘客、行程） | 业务读写 |
| `cs_*` | 客服模块表（会话、工单、FAQ、消息、通知） | 业务读写 |
| `dict_*` | 数据字典表 | 平台维护 |

### 1.3 建模示例

```sql
CREATE TABLE IF NOT EXISTS cs_tickets (
    id           BIGINT       PRIMARY KEY,
    created_at   TIMESTAMPTZ(6),
    updated_at   TIMESTAMPTZ(6),
    deleted_at   TIMESTAMPTZ(6),
    created_by   BIGINT       NOT NULL DEFAULT 0,
    updated_by   BIGINT       NOT NULL DEFAULT 0,
    ticket_no    VARCHAR(50)  NOT NULL DEFAULT '',
    passenger_id BIGINT       NOT NULL DEFAULT 0,
    order_id     BIGINT,
    status       BIGINT       NOT NULL DEFAULT 0
);
COMMENT ON TABLE cs_tickets IS '客服工单';
CREATE UNIQUE INDEX IF NOT EXISTS uk_cs_tickets_ticket_no ON cs_tickets (ticket_no);
CREATE INDEX IF NOT EXISTS idx_cs_tickets_deleted_at ON cs_tickets (deleted_at);
```

### 1.4 SQL 红线

1. **新增表/字段/索引/约束前先评审**，通过后**只通过迁移脚本**落地；**禁止 `AutoMigrate`**（项目 `initialize/db.go` 明确不做迁移）。
2. `UPDATE/DELETE` 必须带 WHERE，禁止全表更新/删除。
3. 必须参数绑定，禁止字符串拼接 SQL。
4. 禁止在本地库冗余外部服务业务明细（只存 `order_id` 等关联 ID）。
5. 联表/聚合尽量在 `service` 层用 GORM 链式调用，复杂查询走脚本评审。

---

## 2. pgvector 向量库规范

### 2.1 基础约定

- 向量能力由**同一 PostgreSQL 实例**的 `pgvector` 扩展提供，禁止引入独立向量库。
- **范围限定**：本次不新增资料外的向量业务；是否使用向量能力，以现有 FAQ/知识库实现为准。
- 启用：`CREATE EXTENSION IF NOT EXISTS vector;`（属 DDL，需确认后执行）。
- 维度统一管理（如 `vector(1536)`），不同模型维度用不同表或分表，禁止混存。

### 2.2 建表与索引

```sql
CREATE TABLE IF NOT EXISTS ride_faq_embeddings (
    id         BIGINT      PRIMARY KEY,
    biz_id     BIGINT      NOT NULL DEFAULT 0,
    content    TEXT        NOT NULL DEFAULT '',
    embedding  vector(1536),
    created_at TIMESTAMPTZ(6)
);
-- HNSW 索引（余弦距离），百万级以内优先；数据量小可用 ivfflat
CREATE INDEX IF NOT EXISTS idx_faq_embedding_hnsw
    ON ride_faq_embeddings USING hnsw (embedding vector_cosine_ops);
```

### 2.3 查询示例

```sql
SELECT id, biz_id, content, 1 - (embedding <=> $1) AS score
FROM ride_faq_embeddings
ORDER BY embedding <=> $1
LIMIT 5;
```

### 2.4 使用规范

1. 向量写入由**业务侧异步化**（Kafka 消费或任务），禁止在同步接口主链路里做 embedding 计算。
2. 先做业务过滤（租户/状态/分类）再做向量检索，禁止全表向量扫描。
3. 阈值与 topK 走配置，禁止硬编码。
4. Go 侧使用参数绑定传向量（如 `'[0.1,0.2,...]'::vector` 通过占位符传入），禁止拼接。
5. 向量字段不参与 GORM 常规 CRUD 结构体时，单独建查询 DTO，避免大字段拖慢列表查询。

---

## 3. Kafka 规范

### 3.1 Topic 命名

```
radish.<域>.<事件>.<版本>
例：radish.order.created.v1、radish.cs.ticket.closed.v1
```

- 全小写、点分、带版本号；禁止无版本、禁止中文、禁止随意缩写。
- 一个事件一个 topic，禁止多事件复用同一 topic。

**事件范围**：仅用于**已存在的**订单状态、车辆指令、故障、客服事件异步解耦，不新增范围外事件。
**契约冻结**：topic 与事件契约需评审后冻结；冻结后的变更按 §3.4 走版本升级。

### 3.2 生产侧

1. 消息体统一 JSON，包含：`eventId`（雪花或 UUID）、`eventType`、`timestamp`、`traceId`、`data`。
2. 生产者在**事务提交后**发送，禁止事务内发消息（避免回滚后消息已发）。
3. 关键事件开启幂等与 ack 配置；发送失败记录日志并落"待补偿"表或重试队列。
4. 禁止把 Kafka 当同步 RPC 使用。

### 3.3 消费侧

1. 消费者必须**幂等**（按 `eventId` 去重/唯一索引）。
2. 消费失败：有限重试 → 死信 topic（`radish.<域>.<事件>.dlq.v1`）→ 告警 + 人工补偿。
3. 禁止单条失败阻塞整个分区（记录错误继续），避免消费堆积。
4. 消费逻辑写在 `service` 层，禁止写在 `api/v1`。

### 3.4 契约

- 事件 schema 变更必须**向后兼容**（新增可选字段，不改语义），不兼容则升 topic 版本。
- 事件字典与 topic 清单维护在文档中，新增 topic 需在 PR 中说明并同步文档。

---

## 4. DDL 变更流程（强制）

1. **先询问确认**：任何建表/改字段/加索引都要先给出方案并等待确认（AGENT 红线）。
2. **脚本化**：DDL 与回滚脚本成对提交到 `sql/`，命名 `YYYYMMDD_<模块>_<动作>.sql`。
3. **评审**：SQL 随 PR 评审；脚本必须幂等（`IF NOT EXISTS`）。
4. **执行**：由负责人在约定窗口执行，禁止 AI/程序自动连接生产执行。
5. **验证**：执行后校验字段类型、索引、注释与 Go model 一致。
6. **记录**：变更记录写入 PR 描述（不写工单编号）。

---

## 5. 数据访问红线（汇总）

1. 业务读写统一 `global.LUOBO_DB`（luobo 库）；`global.GVA_DB` 仅框架 `sys_*` 系统表，**禁止建业务表、禁止写业务数据**（见 §0）。
2. 禁止 SQL 拼接、禁止无 WHERE 写操作、禁止物理删除业务数据。
3. 禁止把密钥、连接串写进代码或日志。
4. 禁止在接口主链路做全表扫描、深翻页未加索引的排序。
5. 大批量操作必须分页/限流，并提供开关与回滚预案。
