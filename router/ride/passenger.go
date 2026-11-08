package ride

import (
	"github.com/gin-gonic/gin"
)

type PassengerRouter struct{}

func (s *PassengerRouter) InitPassengerRouter(Router *gin.RouterGroup) {
	passengerRouter := Router.Group("ride/passenger")
	{
		passengerRouter.GET("list", passengerApi.ListPassengers)
		passengerRouter.GET("detail", passengerApi.GetPassenger)
	}
}
