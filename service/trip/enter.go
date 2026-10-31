package trip

// ServiceGroup 出行运营服务集合（模块1/2/3/4 后台视角，postgres 库）
// 集成方式：在 server/service/enter.go 的 ServiceGroup 中追加 TripServiceGroup trip.ServiceGroup
type ServiceGroup struct {
	TripService
}

var TripServiceApp = new(TripService)
