package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsTicketLogRouter struct{}

func (r *CsTicketLogRouter) InitCsTicketLogRouter(Router *gin.RouterGroup) {
	ticketLogRouter := Router.Group("ticketLog").Use(middleware.OperationRecord())
	ticketLogRouterWithoutRecord := Router.Group("ticketLog")
	{
		ticketLogRouter.POST("", csTicketLogApi.CreateCsTicketLog) // 新增工单流水
	}
	{
		ticketLogRouterWithoutRecord.GET("", csTicketLogApi.GetCsTicketLog)         // 获取工单流水详情
		ticketLogRouterWithoutRecord.GET("list", csTicketLogApi.GetCsTicketLogList) // 分页获取工单流水列表
	}
}
