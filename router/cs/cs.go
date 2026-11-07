package cs

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsRouter struct{}

// InitCsRouter 注册客服中心路由
// 路由前缀：/cs
func (s *CsRouter) InitCsRouter(Router *gin.RouterGroup) {
	router := Router.Group("cs").Use(middleware.OperationRecord())
	readonly := Router.Group("cs")
	api := v1.ApiGroupApp.CsApiGroup.CsApi
	{
		router.POST("ticket", api.CreateTicket)
		router.PUT("ticket", api.UpdateTicket)
		router.PUT("ticket/status", api.UpdateTicketStatus)

		router.POST("faqCategory", api.CreateFaqCategory)
		router.PUT("faqCategory", api.UpdateFaqCategory)
		router.DELETE("faqCategory/:id", api.DeleteFaqCategory)

		router.POST("faq", api.CreateFaq)
		router.PUT("faq", api.UpdateFaq)
		router.DELETE("faq/:id", api.DeleteFaq)
	}
	{
		readonly.GET("session/list", api.GetSessionList)
		readonly.GET("ticket/list", api.GetTicketList)
		readonly.GET("ticket/:id", api.GetTicket)
		readonly.GET("faqCategory/list", api.GetFaqCategoryList)
		readonly.GET("faq/list", api.GetFaqList)
		readonly.GET("notification/list", api.GetNotificationList)
	}
}
