package cs

// RouterGroup 客服中心路由集合
// 集成方式：在 server/router/enter.go 的 RouterGroup 中追加 Cs cs.RouterGroup
type RouterGroup struct {
	CsRouter
}
