package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsTicketRouter struct{}

// InitCsTicketRouter 注册客服工单路由
// 全部工单接口（含创建工单）均注册在 PrivateGroup，统一走 JWT + Casbin 鉴权
func (r *CsTicketRouter) InitCsTicketRouter(Router *gin.RouterGroup) {
	ticketRouter := Router.Group("ticket").Use(middleware.OperationRecord())
	ticketRouterWithoutRecord := Router.Group("ticket")
	{
		ticketRouter.POST("", csTicketApi.CreateCsTicket)                 // 创建客服工单
		ticketRouter.PUT("", csTicketApi.UpdateCsTicket)                  // 修改客服工单
		ticketRouter.PUT("status", csTicketApi.UpdateCsTicketStatus)      // 修改客服工单状态
		ticketRouter.DELETE(":id", csTicketApi.DeleteCsTicket)            // 删除客服工单
	}
	{
		ticketRouterWithoutRecord.GET(":id", csTicketApi.GetCsTicket)         // 获取客服工单详情
		ticketRouterWithoutRecord.GET("list", csTicketApi.GetCsTicketList)    // 分页获取客服工单列表
	}
}
