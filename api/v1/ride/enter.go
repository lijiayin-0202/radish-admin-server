package ride

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	OrderApi
	PassengerApi
	VehicleApi
	DispatchApi
}

var orderService = service.ServiceGroupApp.RideServiceGroup.OrderService
var passengerService = service.ServiceGroupApp.RideServiceGroup.PassengerService
var vehicleService = service.ServiceGroupApp.RideServiceGroup.VehicleService
var dispatchService = service.ServiceGroupApp.RideServiceGroup.DispatchService
