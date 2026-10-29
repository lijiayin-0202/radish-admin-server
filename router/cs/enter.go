package cs

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

// RouterGroup 客服中心与客服运营端路由组
// 统一挂载在 /api/cs 前缀下，由 initialize/router_biz.go 注册
type RouterGroup struct {
	CsFaqRouter
	CsTicketRouter
	CsSessionRouter
	CsMessageRouter
}

var (
	csFaqApi     = api.ApiGroupApp.CsApiGroup.CsFaqApi
	csTicketApi  = api.ApiGroupApp.CsApiGroup.CsTicketApi
	csSessionApi = api.ApiGroupApp.CsApiGroup.CsSessionApi
	csMessageApi = api.ApiGroupApp.CsApiGroup.CsMessageApi
)
