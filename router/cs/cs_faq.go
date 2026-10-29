package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsFaqRouter struct{}

func (r *CsFaqRouter) InitCsFaqRouter(Router *gin.RouterGroup) {
	faqRouter := Router.Group("faq").Use(middleware.OperationRecord())
	faqRouterWithoutRecord := Router.Group("faq")
	{
		faqRouter.POST("", csFaqApi.CreateCsFaq)   // 创建FAQ问答
		faqRouter.PUT("", csFaqApi.UpdateCsFaq)    // 更新FAQ问答
		faqRouter.DELETE("", csFaqApi.DeleteCsFaq) // 删除FAQ问答
	}
	{
		faqRouterWithoutRecord.GET("", csFaqApi.GetCsFaq)         // 获取FAQ问答详情
		faqRouterWithoutRecord.GET("list", csFaqApi.GetCsFaqList) // 分页获取FAQ问答列表
	}
}
