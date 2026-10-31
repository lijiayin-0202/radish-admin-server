package platform

// ServiceGroup 平台配置、权限与统计大屏服务集合
// 集成方式：在 server/service/enter.go 的 ServiceGroup 中追加 PlatformServiceGroup platform.ServiceGroup
type ServiceGroup struct {
	PlatformService
}

var PlatformServiceApp = new(PlatformService)
