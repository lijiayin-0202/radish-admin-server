package ride

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (s *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("trip/order")
	orderRouterWithRecord := Router.Group("trip/order").Use(middleware.OperationRecord())
	{
		orderRouter.GET("list", orderApi.ListOrders)
		orderRouter.GET("detail", orderApi.GetOrder)
	}
	{
		orderRouterWithRecord.POST("dispatch", orderApi.Dispatch)
		orderRouterWithRecord.POST("advance", orderApi.Advance)
	}
}
