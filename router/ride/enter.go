package ride

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	OrderRouter
	PassengerRouter
	VehicleRouter
	DispatchRouter
}

var orderApi = v1.ApiGroupApp.RideApiGroup.OrderApi
var passengerApi = v1.ApiGroupApp.RideApiGroup.PassengerApi
var vehicleApi = v1.ApiGroupApp.RideApiGroup.VehicleApi
var dispatchApi = v1.ApiGroupApp.RideApiGroup.DispatchApi
