package trip

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TripRouter struct{}

// InitTripRouter 注册出行运营路由
// 路由前缀：/trip
func (s *TripRouter) InitTripRouter(Router *gin.RouterGroup) {
	router := Router.Group("trip").Use(middleware.OperationRecord())
	readonly := Router.Group("trip")
	api := v1.ApiGroupApp.TripApiGroup.TripApi
	{
		router.POST("passenger", api.CreatePassenger)
		router.PUT("passenger", api.UpdatePassenger)
		router.DELETE("passenger/:id", api.DeletePassenger)

		router.POST("area", api.CreateArea)
		router.PUT("area", api.UpdateArea)
		router.DELETE("area/:id", api.DeleteArea)

		router.POST("station", api.CreateStation)
		router.PUT("station", api.UpdateStation)
		router.DELETE("station/:id", api.DeleteStation)

		router.POST("vehicle", api.CreateVehicle)
		router.PUT("vehicle", api.UpdateVehicle)
		router.DELETE("vehicle/:id", api.DeleteVehicle)

		router.POST("command", api.CreateCommand)
		router.PUT("command", api.UpdateCommand)

		router.POST("fault", api.CreateFault)
		router.PUT("fault", api.UpdateFault)

		router.PUT("incident", api.UpdateIncident)
	}
	{
		readonly.GET("passenger/list", api.GetPassengerList)
		readonly.GET("order/list", api.GetOrderList)
		readonly.GET("order/statusLogs", api.GetOrderStatusLogs)
		readonly.GET("area/list", api.GetAreaList)
		readonly.GET("station/list", api.GetStationList)
		readonly.GET("vehicle/list", api.GetVehicleList)
		readonly.GET("dispatch/list", api.GetDispatchList)
		readonly.GET("command/list", api.GetCommandList)
		readonly.GET("fault/list", api.GetFaultList)
		readonly.GET("boarding/list", api.GetBoardingList)
		readonly.GET("incident/list", api.GetIncidentList)
	}
}
