# 分册 2｜技术栈与工程结构

> 归属总纲：`../1-总纲/AGENT.md`
> 适用范围：`radish-admin-server`（gin-vue-admin 后台单体服务，Gin + GORM）
> 范围边界：后台前端页面沿用 gin-vue-admin 现成实现，不在本次后端任务中重复开发。

---

## 1. 选型口径（唯一标准答案）

| 维度 | 选型 | 版本/说明 | 禁止事项 |
|------|------|-----------|----------|
| 管理后台框架 | **gin-vue-admin（GVA）** | Server 端 Go 模块 `github.com/flipped-aurora/gin-vue-admin/server`，Go 1.24+ | 禁止替换路由/鉴权底盘，禁止手改 `core/`、`middleware/` |
| 关系数据库 | **PostgreSQL（gva / luobo 双库）** | gva 库 → `global.GVA_DB`（仅框架系统表，业务只读、禁建业务表）；luobo 库 → `global.LUOBO_DB`（业务库，全部业务读写） | 禁止 `AutoMigrate`，表结构只通过迁移脚本落地；禁止业务表建在 gva 库 |
| 向量数据库 | **PostgreSQL + pgvector** | 同 PG 实例 `CREATE EXTENSION vector`，列类型 `vector(dim)`；本次不新增资料外的向量业务，是否使用以现有 FAQ/知识库实现为准 | 禁止引入独立向量库（Milvus 等），禁止 float 数组代替 |
| 消息队列 | **Kafka** | 仅用于**已存在的**订单状态、车辆指令、故障、客服事件异步解耦；topic/事件契约评审后冻结 | 禁止用 Redis/DB 轮询替代消息通道，禁止新增范围外事件 |
| 缓存 | **Redis** | `global.GVA_REDIS`，沿用项目既有封装，缓存与会话 | 禁止把 Redis 当持久化存储，禁止另起客户端 |
| 接口文档 | **Swagger**（`docs/`） | 接口注释即文档，变更后 `swag init` | 禁止口头交付不更新文档 |

---

## 2. GVA 后台工程结构（本仓库）

```
radish-admin-server/
├─ api/v1/            # handler：参数校验 + 调 service + 封装返回（按业务分子包）
├─ service/           # 业务层：逻辑、事务编排（按业务分子包 + enter.go 注册）
├─ model/             # 数据模型（结构体 + TableName + gorm字段tag，只承载数据，不写CRUD数据库操作）
│  ├─ common/request  # PageInfo、GetById 等公共入参
│  ├─ common/response # Response、Ok*、Fail* 统一返回
│  └─ <业务>/request、<业务>/response
├─ router/            # 路由注册（enter.go 聚合 + 各模块 Router）
├─ initialize/        # 启动初始化：db、redis、router、snowflake 等
├─ middleware/        # 中间件（鉴权/操作日志/限流…）
├─ global/            # 全局变量：GVA_DB、LUOBO_DB、GVA_LOG、GVA_CONFIG、GVA_REDIS、Snow
├─ config/            # 配置结构体与 config*.yaml
├─ utils/             # 公共工具（校验、JWT claims、上传…）
├─ resource/template  # 代码生成模板（autocode）
├─ docs/              # swagger 文档产物
└─ sql/               # 建表/迁移脚本（脚本化，禁止随意执行）
```

### 全局变量使用口径

| 变量 | 用途 | 备注 |
|------|------|------|
| `global.GVA_DB` | **gva 库**：gin-vue-admin 原生系统表（`sys_*`、文件、jwt、通知、字典、权限等） | 业务代码只读，**禁止建业务表、禁止写业务数据** |
| `global.LUOBO_DB` | **luobo 库**：业务库，全部业务表与业务 CRUD 走它 | 业务 model 统一用它；禁止 `gorm.Open` 另建连接、禁止手写 DSN |
| `global.GVA_LOG` | zap 日志 | 禁止打印手机号、密码、token |
| `global.Snow` | 雪花 ID（由 `initialize.InitSnowflake` 初始化） | 主键、业务单号统一用它 |
| `global.GVA_CONFIG` | 配置（含 `LuoboPgsql`） | 敏感值走环境变量 |

---

## 3. 跨服务调用（仅出站：GVA 主动发起 HTTP/gRPC）

本仓库是 **gin-vue-admin 单体服务**，全部业务在本仓库内按 GVA 原生分层（`api/v1` / `service` / `model` / `router`）实现。需要外部能力时，只允许 **GVA 主动发起 HTTP/gRPC 调用外部服务**。

```
本仓库内的调用落点（不新增目录）：
service/<业务>/<外部服务>_client.go   # 出站 client：地址、超时、重试、错误转换
service/<业务>/<业务>.go              # 业务编排：调 client + 本地库读写
```

**调用约定**
1. 只允许出站：GVA → 外部服务（HTTP/gRPC）；外部服务不反向依赖本仓库内部包。
2. client 统一放在 `service/` 下，**`api/v1` 不得直连外部服务**，禁止在 handler 里写 `http.Client` / `grpc.Dial`。
3. 外部服务地址、超时、重试次数、开关全部走配置（`config/config.go` + yaml/环境变量，属 Ask），**禁止硬编码 URL**。
4. 必须设置超时（默认 3~5s）与有限重试（≤2 次，且仅幂等请求），禁止无限重试造成雪崩。
5. 外部错误必须转换为本仓库统一返回语义，禁止把外部堆栈、原始报文透出给前端。
6. 订单/行程/司机等外部业务明细**一律调外部服务获取**，本仓库只存 `order_id` 等关联 ID，禁止在本地库冗余外部业务明细。
7. 调用需记录 trace、耗时、结果码，敏感字段脱敏。
8. 接口契约以外部服务提供的接口文档为准，本仓库只负责调用与错误转换。

---

## 4. 环境与配置

| 环境 | 配置文件 | 说明 |
|------|----------|------|
| 本地 | `config.dev.yaml` | 小组统一开发配置，本地启动默认读取 |
| 联调/测试 | `config.yaml`（环境注入） | 敏感项走环境变量，禁止硬编码 |
| 生产 | 环境变量/配置中心 | 禁止提交密钥、生产地址到仓库 |

- 新增配置项：先在 `config/config.go` 结构体补字段（**Ask 权限**），再补 yaml，最后补文档。
- 启动顺序相关：`initialize` 中各初始化函数保持幂等，新增初始化函数挂到 `main.go` 的初始化链（**修改 `main.go` 属 Ask**）。

---

## 5. 依赖管理

- 新增 Go 依赖属于 **Ask 权限**：优先复用现有依赖（`go.mod` 已含 gorm、zap、redis、casbin、excelize 等）。
- 执行 `go mod tidy` 前先确认不会删除间接依赖；提交前 `go build ./...` 必须全绿。
- 禁止引入与选型冲突的库（例如另起 ORM、另起 Web 框架、另起 MQ 客户端）。

---

## 6. 工程结构红线

1. 新业务模块只能新增目录/文件，**不得重排既有目录结构**。
2. 框架底盘（`core/`、`middleware/`、`utils/`、`global/`）只在"框架类变更"分支修改，且单独 PR。
3. 业务模块之间不得横向依赖对方的 `api` 层；跨模块复用只能调对方的 `service`。
4. Swagger/接口文档随代码同步更新，禁止"先上线后补文档"。
