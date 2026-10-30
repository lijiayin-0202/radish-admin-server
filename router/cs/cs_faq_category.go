package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsFaqCategoryRouter struct{}

func (r *CsFaqCategoryRouter) InitCsFaqCategoryRouter(Router *gin.RouterGroup) {
	faqCategoryRouter := Router.Group("faqCategory").Use(middleware.OperationRecord())
	faqCategoryRouterWithoutRecord := Router.Group("faqCategory")
	{
		faqCategoryRouter.POST("", csFaqCategoryApi.CreateCsFaqCategory)   // 新增FAQ分类
		faqCategoryRouter.PUT("", csFaqCategoryApi.UpdateCsFaqCategory)    // 修改FAQ分类
		faqCategoryRouter.DELETE("", csFaqCategoryApi.DeleteCsFaqCategory) // 删除FAQ分类
	}
	{
		faqCategoryRouterWithoutRecord.GET("", csFaqCategoryApi.GetCsFaqCategory)         // 获取FAQ分类详情
		faqCategoryRouterWithoutRecord.GET("list", csFaqCategoryApi.GetCsFaqCategoryList) // 分页获取FAQ分类列表
	}
}
