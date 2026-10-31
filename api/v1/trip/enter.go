package trip

// ApiGroup 出行运营接口集合（模块1/2/3/4 后台视角）
// 集成方式：在 server/api/v1/enter.go 的 ApiGroup 中追加 TripApiGroup trip.ApiGroup
type ApiGroup struct {
	TripApi
}
