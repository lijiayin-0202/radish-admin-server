package cs

// ServiceGroup 客服中心服务集合（模块5，luobo 库）
// 集成方式：在 server/service/enter.go 的 ServiceGroup 中追加 CsServiceGroup cs.ServiceGroup
type ServiceGroup struct {
	CsService
}

var CsServiceApp = new(CsService)
