package ride

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	rideModel "github.com/flipped-aurora/gin-vue-admin/server/model/ride"

	"github.com/gin-gonic/gin"
)

type OrderApi struct{}

// ListOrders 后台订单列表
// @Tags     RideOrder
// @Summary  分页查询订单
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    page      query  int     false "页码"
// @Param    pageSize  query  int     false "每页数量"
// @Param    status    query  int     false "订单状态 0-7"
// @Param    keyword   query  string  false "订单号或乘客ID"
// @Success  200 {object} response.Response{data=ride.OrderListResult} "订单列表"
// @Router   /trip/order/list [get]
func (a *OrderApi) ListOrders(c *gin.Context) {
	var query rideModel.OrderListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := orderService.ListOrders(query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetOrder 订单详情
// @Tags     RideOrder
// @Summary  查询订单详情
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    orderId  query  string  true "订单ID，如 ord_12"
// @Success  200 {object} response.Response{data=object} "订单详情"
// @Router   /trip/order/detail [get]
func (a *OrderApi) GetOrder(c *gin.Context) {
	orderID := c.Query("orderId")
	if orderID == "" {
		response.FailWithMessage("orderId 不能为空", c)
		return
	}
	result, err := orderService.GetOrder(orderID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// Dispatch 后台派车
// @Tags     RideOrder
// @Summary  派单并确认车辆
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      ride.DispatchRequest  true "订单ID与可选车辆ID"
// @Success  200 {object} response.Response{data=object} "派车结果"
// @Router   /trip/order/dispatch [post]
func (a *OrderApi) Dispatch(c *gin.Context) {
	var payload rideModel.DispatchRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := orderService.Dispatch(payload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "派车成功", c)
}

// Advance 推进订单状态
// @Tags     RideOrder
// @Summary  后台推进订单状态（到达/开始/完成/异常关闭）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      ride.CommandRequest  true "订单ID与命令类型"
// @Success  200 {object} response.Response{data=object} "最新订单"
// @Router   /trip/order/advance [post]
func (a *OrderApi) Advance(c *gin.Context) {
	var payload rideModel.CommandRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := orderService.Advance(payload)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "操作成功", c)
}
