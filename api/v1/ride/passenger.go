package ride

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	rideModel "github.com/flipped-aurora/gin-vue-admin/server/model/ride"

	"github.com/gin-gonic/gin"
)

type PassengerApi struct{}

// ListPassengers 后台乘客列表
// @Tags     RidePassenger
// @Summary  分页查询乘客
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    page            query  int     false "页码"
// @Param    pageSize        query  int     false "每页数量"
// @Param    keyword         query  string  false "乘客号/手机号/昵称/乘客ID"
// @Param    realNameStatus  query  int     false "实名状态"
// @Param    accountStatus   query  int     false "账户状态"
// @Success  200 {object} response.Response{data=ride.PassengerListResult} "乘客列表"
// @Router   /trip/passenger/list [get]
func (a *PassengerApi) ListPassengers(c *gin.Context) {
	var query rideModel.PassengerListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := passengerService.ListPassengers(query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetPassenger 后台乘客详情
// @Tags     RidePassenger
// @Summary  查询乘客详情
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    passengerId  query  string  true "乘客ID"
// @Success  200 {object} response.Response{data=ride.PassengerDetailResult} "乘客详情"
// @Router   /trip/passenger/detail [get]
func (a *PassengerApi) GetPassenger(c *gin.Context) {
	passengerID := c.Query("passengerId")
	if passengerID == "" {
		response.FailWithMessage("passengerId 不能为空", c)
		return
	}
	result, err := passengerService.GetPassenger(passengerID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
