# 模块 3 运营区域与站点可视化管理设计

## 目标

在现有后台管理系统中提供“运营调度管理”页面，供已登录的后台人员完成运营区域和站点的查询、新增、编辑、启停及运营时段维护。页面只通过后台管理服务访问模块 3 微服务，不直接暴露微服务地址给浏览器。

## 范围与边界

### 本次实现

- 运营区域：分页查询、按城市/状态筛选、新增、编辑、启用/临时关闭/停运、运营时段维护。
- 站点：按运营区域筛选、名称筛选、新增、编辑、启用/停用。
- 后台管理端作为 BFF（Backend For Frontend）：鉴权后转发模块 3 HTTP 请求，并转换为当前后台的统一响应格式。
- 前端在 `src/view/module3/` 与 `src/api/module3.js` 中新增独立业务文件。

### 明确不做

- 不修改模块 3 Proto、数据库表、迁移文件、派单逻辑或雪花 ID 规则。
- 不修改前端 `src/core/`、后台 `core/`、JWT/Casbin 中间件。
- 不通过 SQL 自动插入菜单、角色或权限数据；页面完成后由管理员在现有“菜单管理”界面登记菜单。
- 不新增 npm 或 Go 依赖。

## 架构与数据流

```text
Vue 页面
  -> /api/module3/*（携带现有 x-token）
  -> radish-admin-server 私有路由（JWT + Casbin）
  -> 模块 3 HTTP 服务（http://127.0.0.1:8000）
  -> PostgreSQL
```

前端只识别后台管理服务返回的 `{ code, data, msg }`；后台适配层负责处理模块 3 的原始 JSON、超时和网络错误。微服务不可用时，后台返回已有统一失败结构，不返回底层地址、令牌或内部错误详情。

## 配置

后台管理端新增配置项：

```yaml
ride-hailing:
  base-url: http://127.0.0.1:8000
  request-timeout: 3s
```

- `base-url` 仅用于后台管理服务调用模块 3。
- 实际部署地址由环境配置覆盖；前端代码不保存此地址。
- 超时后返回统一失败响应，且不执行重试写操作，避免新增或更新被重复提交。

## 后台管理端接口

所有接口注册在现有私有路由组中，继续使用项目 JWT 和 Casbin 校验。

| 后台管理端接口 | 方法 | 对应模块 3 接口 | 用途 |
| --- | --- | --- | --- |
| `/module3/service-area/list` | POST | `GET /v1/admin/service-areas` | 区域分页查询 |
| `/module3/service-area/detail` | POST | `GET /v1/admin/service-areas/{id}` | 读取区域和当前完整运营时段，用于编辑回填 |
| `/module3/service-area/create` | POST | `POST /v1/admin/service-areas` | 新增区域 |
| `/module3/service-area/update` | POST | `PATCH /v1/admin/service-areas/{id}` | 编辑区域 |
| `/module3/service-area/change-status` | POST | `PATCH /v1/admin/service-areas/{id}/status` | 区域启停 |
| `/module3/service-area/schedules` | POST | `PUT /v1/admin/service-areas/{id}/schedules` | 保存运营时段 |
| `/module3/station/list` | POST | `GET /v1/admin/stations` | 站点查询 |
| `/module3/station/create` | POST | `POST /v1/admin/stations` | 新增站点 |
| `/module3/station/update` | POST | `PATCH /v1/admin/stations/{id}` | 编辑站点 |
| `/module3/station/change-status` | POST | `PATCH /v1/admin/stations/{id}/status` | 站点启停 |

## 字段约束

- 区域状态只使用既有值：`0 临时关闭`、`1 运营中`、`2 停运`。
- 区域编码 `areaCode`、站点编码 `stationCode` 和站点所属区域 `serviceAreaId` 仅在新建时填写；编辑时不提供这些字段。
- 运营时段沿用微服务字段：星期、特殊日期、开始分钟、结束分钟、是否启用、公告；页面以 `HH:mm` 显示，提交时转换为 `0-1439` 分钟。
- 时段编辑必须先读取区域详情，再执行覆盖保存；禁止用空列表直接覆盖未知的既有时段。
- 站点列表的微服务协议没有分页字段，因此后台保留原协议；前端只对当前筛选结果做本地表格分页，并明确显示“当前运营区域内站点”。

## 前端页面

入口页面：`src/view/module3/operationManagement.vue`。

- 使用 Element Plus 的标签页：`运营区域` 与 `站点管理`。
- 区域页：顶部筛选表单、列表表格、状态标签、添加/编辑弹窗、运营时段弹窗。
- 站点页：运营区域下拉框为必选筛选条件，列表表格、添加/编辑弹窗、启停确认框。
- 写操作完成后只刷新对应列表；失败时沿用项目全局消息提示。
- 新菜单组件路径建议填写为 `view/module3/operationManagement.vue`，路由名使用 `Module3OperationManagement`。

## 菜单与权限登记

页面代码不直接改系统菜单表。管理员在现有“超级管理员 -> 菜单管理”中新增：

- 标题：`运营调度管理`
- 路由：`module3OperationManagement`
- 组件：`view/module3/operationManagement.vue`
- 路由名称：`Module3OperationManagement`

随后在现有角色菜单关联页面分配相应角色。这样不会冲突或覆盖团队已有菜单初始化数据。

## 验收与测试

1. 后台管理端：对模块 3 HTTP 客户端使用 `httptest` 验证 URL、方法、请求体、超时和错误转换；执行 `go test ./...` 与 `go build ./...`。
2. 前端：执行 `npm run build`；人工验证区域和站点的新增、编辑、启停、筛选与运营时段保存。
3. 联调：启动后台管理端和模块 3 微服务，使用登录后的后台页面完成一条区域、站点的全链路操作；记录在 `development` 目录的任务验收记录中。

## 风险与处理

- 模块 3 HTTP 接口当前不包含后台 JWT 校验：仅允许后台管理端服务端调用，部署时应通过内网或网关限制微服务管理端口。
- 动态菜单依赖后台现有系统菜单数据：本次不自动写入，避免越过团队权限和菜单管理流程。
- 多人同时修改前端动态菜单、后台 `initialize/router_biz.go` 或配置样例时会有冲突：模块业务文件保持独立，公共入口由合并时小范围处理。
