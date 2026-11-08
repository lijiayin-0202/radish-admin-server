package ride

import (
	"github.com/gin-gonic/gin"
)

type VehicleRouter struct{}

func (s *VehicleRouter) InitVehicleRouter(Router *gin.RouterGroup) {
	vehicleRouter := Router.Group("ride/vehicle")
	{
		vehicleRouter.GET("list", vehicleApi.ListVehicles)
		vehicleRouter.GET("detail", vehicleApi.GetVehicle)
	}
}
