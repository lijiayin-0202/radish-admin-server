package platform

// RouterGroup 平台配置、权限与统计大屏路由集合
// 集成方式：在 server/router/enter.go 的 RouterGroup 中追加 Platform platform.RouterGroup
type RouterGroup struct {
	PlatformRouter
}
