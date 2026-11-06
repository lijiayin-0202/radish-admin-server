package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsNotificationRouter struct{}

func (r *CsNotificationRouter) InitCsNotificationRouter(Router *gin.RouterGroup) {
	notificationRouter := Router.Group("notification").Use(middleware.OperationRecord())
	notificationRouterWithoutRecord := Router.Group("notification")
	{
		notificationRouter.PUT("read", csNotificationApi.ReadCsNotification) // 标记通知已读
	}
	{
		notificationRouterWithoutRecord.GET("", csNotificationApi.GetCsNotification)         // 获取通知详情
		notificationRouterWithoutRecord.GET("list", csNotificationApi.GetCsNotificationList) // 分页获取我的通知列表
	}
}
