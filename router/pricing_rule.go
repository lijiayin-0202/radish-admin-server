package router

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// PricingRuleRouter 计价规则路由
type PricingRuleRouter struct{}

// InitPricingRuleRouter 注册计价规则路由（前缀 /pricingRule）
func (s *PricingRuleRouter) InitPricingRuleRouter(Router *gin.RouterGroup) {
	api := &v1.ApiGroupApp.PricingRuleApi
	write := Router.Group("pricingRule").Use(middleware.OperationRecord())
	read := Router.Group("pricingRule")
	{
		write.POST("create", api.CreatePricingRule)
		write.PUT("update", api.UpdatePricingRule)
		write.DELETE("delete", api.DeletePricingRule)
		write.DELETE("deleteByIds", api.DeletePricingRulesByIds)
	}
	{
		read.GET("find", api.FindPricingRule)
		read.GET("list", api.GetPricingRuleList)
	}
}
