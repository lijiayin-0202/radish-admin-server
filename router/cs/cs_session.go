package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsSessionRouter struct{}

func (r *CsSessionRouter) InitCsSessionRouter(Router *gin.RouterGroup) {
	sessionRouter := Router.Group("session").Use(middleware.OperationRecord())
	sessionRouterWithoutRecord := Router.Group("session")
	{
		sessionRouter.POST("", csSessionApi.CreateCsSession)    // 创建客服会话
		sessionRouter.PUT("", csSessionApi.UpdateCsSession)     // 更新客服会话
		sessionRouter.DELETE("", csSessionApi.DeleteCsSession)  // 删除客服会话
	}
	{
		sessionRouterWithoutRecord.GET("", csSessionApi.GetCsSession)         // 获取客服会话详情
		sessionRouterWithoutRecord.GET("list", csSessionApi.GetCsSessionList) // 分页获取客服会话列表
	}
}
