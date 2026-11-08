# 模块 3 BFF 运行时修复记录

## 任务状态

测试通过，待页面验收。

## 现象

运营调度管理页面加载时提示“模块3服务暂不可用”。

## 根因

运行中的 `server.exe` 构建时间早于模块 3 响应兼容代码。模块 3 微服务以 Proto JSON 返回运营区域状态枚举文本，例如：

```text
SERVICE_AREA_STATUS_OPERATING
```

旧二进制无法解析该字段，BFF 将解码失败统一转换为“模块3服务暂不可用”。

## 修复

结束旧 `server.exe` 后，使用当前源码和 `config.dev.yaml` 启动后台服务：

```powershell
go run main.go -c config.dev.yaml
```

当前源码包含 `ServiceAreaStatus.UnmarshalJSON`，可同时兼容 Proto 枚举文本和数值状态。

## 验证

- `go test ./service/module3 ./api/v1/module3 -count=1` 通过。
- 当前源码 BFF 客户端直连模块 3 成功：`areas=10`，`total=42`。
- 后台健康接口：`GET /health` 返回 `200`、`ok`。
- 模块 3 上游接口：`GET /v1/admin/service-areas?page=1&page_size=10` 返回 `200`、`areas=10`、`total=42`。
- 启动日志确认全部 `/module3/*` BFF 路由已注册。

## 页面验收

1. 在浏览器执行 `Ctrl + F5`。
2. 打开“出行运营 → 运营调度管理”。
3. 确认运营区域列表正常显示，不再出现“模块3服务暂不可用”。
4. 切换到“站点管理”，确认可选择运营区域。
