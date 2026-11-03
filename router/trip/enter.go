package trip

// RouterGroup 出行运营路由集合
// 集成方式：在 server/router/enter.go 的 RouterGroup 中追加 Trip trip.RouterGroup
type RouterGroup struct {
	TripRouter
}
