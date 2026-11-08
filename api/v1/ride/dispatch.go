package ride

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	rideModel "github.com/flipped-aurora/gin-vue-admin/server/model/ride"

	"github.com/gin-gonic/gin"
)

type DispatchApi struct{}

// ListDispatchRecords 后台派单记录列表
// @Tags     RideDispatch
// @Summary  分页查询派单记录
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    page      query  int     false "页码"
// @Param    pageSize  query  int     false "每页数量"
// @Param    keyword   query  string  false "派单号/订单号/车牌/订单ID"
// @Param    status    query  int     false "状态 0已创建 1已下发 2车辆确认 3失败 4超时"
// @Success  200 {object} response.Response{data=ride.DispatchRecordListResult} "派单记录列表"
// @Router   /trip/dispatch/list [get]
func (a *DispatchApi) ListDispatchRecords(c *gin.Context) {
	var query rideModel.DispatchRecordListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := dispatchService.ListDispatchRecords(query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}
