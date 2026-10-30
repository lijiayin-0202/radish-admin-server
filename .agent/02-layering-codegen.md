# 分册 3｜分层编码与代码生成规范

> 归属总纲：`../1-总纲/AGENT.md`
> 核心：`api/v1(handler) → service → model`，骨架优先生成，人工只补业务

---

## 1. 分层职责

| 层 | 职责 | 允许 | 禁止 |
|----|------|------|------|
| `router` | 路由注册、中间件挂载 | `Router.Group`、挂 `middleware.OperationRecord()` | 写业务逻辑 |
| `api/v1` | 参数绑定 + `utils.Verify` 校验 + 调 service + `response.*` 返回 | 读 `utils.GetUserID(c)` | 写 SQL、直接操作 DB、跨层调别的 api |
| `service` | 业务逻辑、事务编排、外部服务调用（HTTP/gRPC）、Kafka 收发；**GORM 数据库操作统一写在此** | 事务、日志、`global.LUOBO_DB` 上的链式 CRUD | 写 SQL、直接写 HTTP 响应 |
| `model` | 表结构映射（结构体 + `TableName()` + 字段 tag），**只承载数据** | `global.GVA_MODEL`、`TableName()`、字段 tag | 定义 CRUD 成员方法、业务逻辑、事务编排 |
| `request` | 入参（Search/CreateReq/UpdateReq） | 内嵌 `commonReq.PageInfo` | 与 model 混用 |
| `response` | 出参 | 组合 model 返回 | 暴露敏感字段（密码、token） |

**分层红线**：`api/v1` handler 不写 SQL；业务逻辑与 GORM 数据库操作只在 `service`；`model` 只定义结构与表名，不写 CRUD 方法。
**禁止引入** Kratos、proto、biz、data、cmd、wire 相关代码或目录。

### 1.1 双库隔离（P0，优先于本分册其他分层规则）

| 库 | 句柄 | 存放内容 | 禁止 |
|----|------|----------|------|
| gva 库 | `global.GVA_DB` | gin-vue-admin 原生系统表（`sys_*`、文件、jwt、通知、字典、权限等框架内置表） | **禁止新建任何业务表、业务 model**；业务代码不得往 gva 库写业务数据 |
| luobo 库 | `global.LUOBO_DB` | 全部业务表（`cs_*` 客服、乘客、订单、钱包、优惠券等，含后续新增业务表） | 禁止把业务表建在 gva 库；禁止业务 model 使用 `GVA_DB` |

1. **严格隔离**：`GVA_DB` 只负责框架自带系统功能；所有业务 CRUD、业务表、业务 model 一律使用 `LUOBO_DB`；业务表禁止创建在 gva 库。
2. **`GVA_MODEL` 只是字段结构体（重要）**：业务 model **可以**嵌入 `global.GVA_MODEL`，但**仅为复用 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt` 四个字段**；`GVA_MODEL` **不会自动绑定任何数据库连接**。业务增删改查必须**显式使用 `global.LUOBO_DB` 操作 model 实例**，**禁止用 model 实例变量直接调用 GORM 方法**（Find / Create / Save / Updates / Delete 等），model 实例只作数据载体。
3. **事务统一在 `service` 层开启**：业务事务使用 `global.LUOBO_DB.Transaction`；**禁止多层事务嵌套**；事务内部禁止发起外部 HTTP 调用、Kafka 同步发送。
4. **只读统计**：允许基于 `global.LUOBO_DB` 开 GORM 只读 `Session`；**禁止**代码内单独 `gorm.Open` 新建 DB 连接、禁止手写 DSN。
5. **向量查询参数化**：pgvector 查询必须参数绑定，禁止字符串拼接向量数组（详见分册 7 §2.4）。
6. **SQL 日志**：开发环境可开 GORM `Debug`；测试、生产环境**禁止打印 SQL 绑定参数**。
7. **链式起点必须是句柄**：所有 GORM 链式调用的**第一个接收者**必须是 `global.LUOBO_DB`（业务）或 `global.GVA_DB`（框架 `sys_*` 系统表）；**禁止以 model 结构体实例作为链式起点**。
8. **禁止在 model 上写 CRUD 方法**：不得为业务 model 结构体定义 Find / Create / Update / Delete / Save 等成员方法；数据库操作代码统一放在 `service` 层，`model` 只保留结构体、`TableName()`、字段 tag。
9. **软删除统一走 `Delete`**：用 `global.LUOBO_DB.Delete(&entity)`，由 GORM 自动维护 `DeletedAt`；**禁止手动 `Update` 修改 `deleted_at` 字段**；查询已删除数据必须显式调用 `Unscoped()`。

---

## 2. 命名规范

| 对象 | 规范 | 示例 |
|------|------|------|
| 文件 | 小写 + 下划线，与表名一致 | `cs_tickets.go`、`cs_faq.go` |
| model 结构体 | 大驼峰单数 | `CsTicket` |
| service 结构体 | `<Model>Service` | `CsTicketService` |
| api 结构体 | `<Model>Api` | `CsTicketApi` |
| router 结构体 | `<Model>Router` | `CsTicketRouter` |
| 方法 | `Create/Update/Delete/Get/GetList` + 实体名 | `CreateCsTicket` |
| 入参 | `<Model>Search`、`CreateReq`、`UpdateReq` | `CsTicketCreateReq` |
| 出参 | `<Model>Response`、列表 `CsListResponse` | `CsTicketResponse` |
| 注册文件 | 各层 `enter.go` 聚合 | `service/cs/enter.go` |

**注册三件套（新增模块必须同步）**
- `api/v1/enter.go`：`ApiGroup` 增加 `XxxApiGroup`
- `service/enter.go`：`ServiceGroup` 增加 `XxxServiceGroup`
- `router/enter.go`：`RouterGroup` 增加 `Xxx`
- `initialize/router_biz.go`：调用 `InitXxxRouter`，统一前缀（如 `/api/cs`）

---

## 3. model 规范

```go
type CsTicket struct {
	global.GVA_MODEL                 // 仅复用 ID/CreatedAt/UpdatedAt/DeletedAt 字段（不绑定任何 DB 连接），禁止自造删除字段
	CreatedBy   int64      `json:"createdBy" gorm:"column:created_by;not null;default:0;comment:创建人"`
	UpdatedBy   int64      `json:"updatedBy" gorm:"column:updated_by;not null;default:0;comment:更新人"`
	TicketNo    string     `json:"ticketNo" gorm:"column:ticket_no;type:varchar(50);not null;default:'';comment:工单编号"`
	PassengerID int64      `json:"passengerId" gorm:"column:passenger_id;type:bigint;not null;default:0;comment:乘客ID"`
	OrderID     *int64     `json:"orderId" gorm:"column:order_id;type:bigint;comment:关联订单ID"`
	ClosedAt    *time.Time `json:"closedAt" gorm:"column:closed_at;comment:关闭时间"`
}

func (CsTicket) TableName() string { return "cs_tickets" }
```

- 主键/单号统一雪花 ID：`global.Snow.NextID()`（`strconv.FormatInt` 转字符串存单号）。
- 金额用 `decimal`，禁止 `float`；可空字段用指针（`*int64`、`*time.Time`）。
- 每张表必须写 `TableName()`；字段必须有 `comment`。
- **业务 model 对应 luobo 库表，统一使用 `global.LUOBO_DB`**；gva 库 `sys_*` 系统表**禁止新建业务 model、禁止扩展字段**。
- **CRUD 入口只能是库句柄（写在 `service` 层）**：统一写 `global.LUOBO_DB.操作(&csModel.CsTicket{})`，例如 `global.LUOBO_DB.Model(&csModel.CsTicket{}).Where("status = ?", status).Find(&list)`、`global.LUOBO_DB.Create(&entity)`；model 实例只承载数据，**不得作为数据库入口**。
- **禁止写法**：用 model 实例变量调用 GORM 的 Find / Create / Save / Updates / Delete 一律禁止（实例上没有连接，即使能编译也不允许）。
- **禁止**自造 `DeletedAt`、`IsDelete` 等删除字段；软删除统一 `global.LUOBO_DB.Delete(&entity)`，由 GORM 自动维护 `DeletedAt`。
- **禁止手动改删除标记**：不得用 `Update` 修改 `deleted_at`；默认查询自动排除已删数据，**查已删除数据必须显式调用 `Unscoped()`**。
- **禁止在 model 上写 CRUD 方法**：不得为业务 model 定义 Find / Create / Update / Delete / Save 等成员方法；GORM 调用代码统一写在 `service` 层。

---

## 4. service 规范

```go
func (s *CsTicketService) CreateCsTicket(operatorID uint, req csReq.CsTicketCreateReq) (err error) {
	entity := csModel.CsTicket{
		TicketNo:    strconv.FormatInt(global.Snow.NextID(), 10), // 后端生成，前端不传
		PassengerID: req.PassengerID,
		CreatedBy:   int64(operatorID),
		UpdatedBy:   int64(operatorID),
	}
	return global.LUOBO_DB.Create(&entity).Error
}
```

- **业务数据源统一 `global.LUOBO_DB`**（luobo 库）；`global.GVA_DB` 仅用于框架 `sys_*` 等系统表，业务 CRUD 禁止使用。禁止另起连接、禁止手写 DSN。
- 列表：`global.LUOBO_DB.Scopes(search.PageInfo.Paginate())` + `Count(&total)`，排序明确（链式首接收者必须是库句柄）。
- 更新：先 `First` 判存在，再赋值，`Save`；**单号、创建人等不可变字段不覆盖**。
- 删除：`global.LUOBO_DB.Delete(&entity)` 软删除，禁止物理删除、禁止全表删除。
- 事务：业务事务统一在 `service` 层用 `global.LUOBO_DB.Transaction(func(tx *gorm.DB) error { ... })` 开启；禁止嵌套事务；事务内禁止外部 HTTP 调用与 Kafka 同步发送；异常必须回滚并记录日志。
- 条件拼接必须用参数绑定 `global.LUOBO_DB.Where("status = ?", status)`，**禁止字符串拼接 SQL**。

---

## 5. api 规范

```go
func (a *CsTicketApi) CreateCsTicket(c *gin.Context) {
	var req csReq.CsTicketCreateReq
	if err := c.ShouldBindJSON(&req); err != nil { response.FailWithMessage(err.Error(), c); return }
	if err := utils.Verify(req, CsTicketCreateVerify); err != nil { response.FailWithMessage(err.Error(), c); return }
	if err := csTicketService.CreateCsTicket(utils.GetUserID(c), req); err != nil {
		global.GVA_LOG.Error("创建客服工单失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c); return
	}
	response.OkWithMessage("创建成功", c)
}
```

- 每个接口必须有 swagger 注释（`@Tags/@Summary/@Param/@Success/@Router`）。
- 校验规则定义在各模块 `enter.go`（如 `CsTicketCreateVerify = utils.Rules{...}`），**不改框架 `utils/validator.go`**。
- 返回统一 `response.OkWithMessage/OkWithDetailed/FailWithMessage`，禁止裸返回 `error` 或自定义结构。
- 分页校验 `utils.Verify(search.PageInfo, utils.PageInfoVerify)`；按 ID 查询用 `request.GetById` + `utils.IdVerify`。

---

## 6. router 规范

```go
func (r *CsTicketRouter) InitCsTicketRouter(Router *gin.RouterGroup) {
	ticketRouter := Router.Group("ticket").Use(middleware.OperationRecord()) // 写操作记录操作日志
	ticketRouterWithoutRecord := Router.Group("ticket")                      // 读操作不记录
	{
		ticketRouter.POST("", csTicketApi.CreateCsTicket)
		ticketRouter.PUT("", csTicketApi.UpdateCsTicket)
		ticketRouter.DELETE("", csTicketApi.DeleteCsTicket)
	}
	{
		ticketRouterWithoutRecord.GET("", csTicketApi.GetCsTicket)
		ticketRouterWithoutRecord.GET("list", csTicketApi.GetCsTicketList)
	}
}
```

- 模块前缀统一（如 `/api/cs`），资源名小写（`faq`、`ticket`、`session`）。
- 写接口必须挂 `middleware.OperationRecord()`；读接口不挂。

---

## 7. 代码生成优先（命令/生成器）

1. **能生成的一律生成**：新建模块的 CRUD 骨架（model / request / response / api / service / router）优先使用 GVA 自带代码生成器（Web「系统工具 → 代码生成器」或对应 autocode 接口/命令）产出，模板来自 `resource/template`。
2. 生成前：先确认表结构（DDL 未确认不生成），字段类型、注释与数据库保持一致。
3. 生成后：
   - 允许：补业务校验规则、补查询条件、补注释、补单号/操作人赋值、补外部服务调用与事务。
   - 禁止：改骨架结构、改命名风格、删 `enter.go` 注册、改返回封装、手改生成模板（`resource/template`）。
4. 生成器覆盖不到的场景（联表聚合、向量检索、Kafka 生产/消费、调用外部服务）由人工在 `service` 层补，写法与本分册保持一致。
5. swagger：`swag init` 重新生成 `docs/`，禁止手改 `docs/docs.go`、`docs/swagger.json`。

---

## 8. 编码红线

1. 禁止跨层调用（api 不碰 DB，router 不写逻辑，service 不写 HTTP 响应）。
2. 禁止硬编码价格、费率、密钥、环境地址；统一走配置或数据库。
3. 禁止全表 `UPDATE/DELETE`；所有写操作带 WHERE。
4. 日志禁止输出手机号、密码、token、身份证等敏感数据。
5. 禁止"顺手重构"：无明确指令不做大范围存量代码改写。
6. 新增依赖、改表、改中间件均需先确认（Ask）。
7. **双库隔离**：业务表、业务 model、业务事务只走 `global.LUOBO_DB`；`global.GVA_DB` 仅框架系统表，禁止建业务表、禁止写业务数据。
