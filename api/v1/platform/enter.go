package platform

// ApiGroup 平台配置、权限与统计大屏接口集合
// 集成方式：在 server/api/v1/enter.go 的 ApiGroup 中追加 PlatformApiGroup platform.ApiGroup
type ApiGroup struct {
	PlatformApi
}
