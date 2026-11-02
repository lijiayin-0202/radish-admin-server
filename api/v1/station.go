package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
)

// StationApi 站点管理接口
type StationApi struct{}

// CreateStation 创建站点
// @Tags     站点管理
// @Summary  创建站点
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.Station true "站点"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /station/create [post]
func (a *StationApi) CreateStation(c *gin.Context) {
	var req tripModel.Station
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.StationService.CreateStation(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateStation 更新站点
// @Tags     站点管理
// @Summary  更新站点
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.Station true "站点"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /station/update [put]
func (a *StationApi) UpdateStation(c *gin.Context) {
	var req tripModel.Station
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.StationService.UpdateStation(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteStation 删除站点
// @Tags     站点管理
// @Summary  删除站点
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.GetById true "主键ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /station/delete [delete]
func (a *StationApi) DeleteStation(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.StationService.DeleteStation(req.Uint()); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteStationsByIds 批量删除站点
// @Tags     站点管理
// @Summary  批量删除站点
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.IdsReq true "id列表"
// @Success  200  {object} response.Response{msg=string} "批量删除成功"
// @Router   /station/deleteByIds [delete]
func (a *StationApi) DeleteStationsByIds(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.StationService.DeleteStationByIds(req.Ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// FindStation 根据id获取站点
// @Tags     站点管理
// @Summary  根据id获取站点
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    request.GetById true "主键ID"
// @Success  200  {object} response.Response{data=tripModel.Station,msg=string} "获取成功"
// @Router   /station/find [get]
func (a *StationApi) FindStation(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := service.ServiceGroupApp.StationService.GetStation(req.Uint())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetStationList 分页获取站点列表
// @Tags     站点管理
// @Summary  分页获取站点列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    tripReq.TripSearch true "分页+筛选"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /station/list [get]
func (a *StationApi) GetStationList(c *gin.Context) {
	var info tripReq.TripSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupApp.StationService.GetStationList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}
