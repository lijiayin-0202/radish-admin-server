package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
)

// ServiceAreaApi 运营区域接口
type ServiceAreaApi struct{}

// CreateServiceArea 创建运营区域
// @Tags     运营区域
// @Summary  创建运营区域
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.ServiceArea true "运营区域"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /serviceArea/create [post]
func (a *ServiceAreaApi) CreateServiceArea(c *gin.Context) {
	var req tripModel.ServiceArea
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.ServiceAreaService.CreateServiceArea(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateServiceArea 更新运营区域
// @Tags     运营区域
// @Summary  更新运营区域
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.ServiceArea true "运营区域"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /serviceArea/update [put]
func (a *ServiceAreaApi) UpdateServiceArea(c *gin.Context) {
	var req tripModel.ServiceArea
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.ServiceAreaService.UpdateServiceArea(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteServiceArea 删除运营区域
// @Tags     运营区域
// @Summary  删除运营区域
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.GetById true "主键ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /serviceArea/delete [delete]
func (a *ServiceAreaApi) DeleteServiceArea(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.ServiceAreaService.DeleteServiceArea(req.Uint()); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteServiceAreasByIds 批量删除运营区域
// @Tags     运营区域
// @Summary  批量删除运营区域
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.IdsReq true "id列表"
// @Success  200  {object} response.Response{msg=string} "批量删除成功"
// @Router   /serviceArea/deleteByIds [delete]
func (a *ServiceAreaApi) DeleteServiceAreasByIds(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.ServiceAreaService.DeleteServiceAreaByIds(req.Ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// FindServiceArea 根据id获取运营区域
// @Tags     运营区域
// @Summary  根据id获取运营区域
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    request.GetById true "主键ID"
// @Success  200  {object} response.Response{data=tripModel.ServiceArea,msg=string} "获取成功"
// @Router   /serviceArea/find [get]
func (a *ServiceAreaApi) FindServiceArea(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := service.ServiceGroupApp.ServiceAreaService.GetServiceArea(req.Uint())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetServiceAreaList 分页获取运营区域列表
// @Tags     运营区域
// @Summary  分页获取运营区域列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    tripReq.TripSearch true "分页+筛选"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /serviceArea/list [get]
func (a *ServiceAreaApi) GetServiceAreaList(c *gin.Context) {
	var info tripReq.TripSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupApp.ServiceAreaService.GetServiceAreaList(info)
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
