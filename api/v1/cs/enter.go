package cs

// ApiGroup 客服中心接口集合（模块5）
// 集成方式：在 server/api/v1/enter.go 的 ApiGroup 中追加 CsApiGroup cs.ApiGroup
type ApiGroup struct {
	CsApi
}
