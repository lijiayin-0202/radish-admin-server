package ride

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	rideModel "github.com/flipped-aurora/gin-vue-admin/server/model/ride"

	"github.com/gin-gonic/gin"
)

type VehicleApi struct{}

// ListVehicles 后台车辆列表
// @Tags     RideVehicle
// @Summary  分页查询车辆
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    page             query  int     false "页码"
// @Param    pageSize         query  int     false "每页数量"
// @Param    keyword          query  string  false "车牌/车辆编号/VIN/车型"
// @Param    operationStatus  query  int     false "运营状态 0停运 1空闲 2接驾中 3待上车 4行程中 5维护"
// @Param    onlineStatus     query  int     false "在线状态 0离线 1在线"
// @Success  200 {object} response.Response{data=ride.VehicleListResult} "车辆列表"
// @Router   /trip/vehicle/list [get]
func (a *VehicleApi) ListVehicles(c *gin.Context) {
	var query rideModel.VehicleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := vehicleService.ListVehicles(query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetVehicle 后台车辆详情
// @Tags     RideVehicle
// @Summary  查询车辆详情
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    vehicleId  query  string  true "车辆ID"
// @Success  200 {object} response.Response{data=ride.VehicleDetailResult} "车辆详情"
// @Router   /trip/vehicle/detail [get]
func (a *VehicleApi) GetVehicle(c *gin.Context) {
	vehicleID := c.Query("vehicleId")
	if vehicleID == "" {
		response.FailWithMessage("vehicleId 不能为空", c)
		return
	}
	result, err := vehicleService.GetVehicle(vehicleID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
