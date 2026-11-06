package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsMessageRouter struct{}

func (r *CsMessageRouter) InitCsMessageRouter(Router *gin.RouterGroup) {
	messageRouter := Router.Group("message").Use(middleware.OperationRecord())
	messageRouterWithoutRecord := Router.Group("message")
	{
		messageRouter.POST("", csMessageApi.CreateCsMessage)   // 保存会话消息
		messageRouter.PUT("", csMessageApi.UpdateCsMessage)    // 更新会话消息
		messageRouter.DELETE("", csMessageApi.DeleteCsMessage) // 删除会话消息
	}
	{
		messageRouterWithoutRecord.GET("", csMessageApi.GetCsMessage)         // 获取会话消息详情
		messageRouterWithoutRecord.GET("list", csMessageApi.GetCsMessageList) // 分页获取会话历史消息（按会话ID查询）
	}
}
