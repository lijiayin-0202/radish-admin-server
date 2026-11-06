package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
)

// OperationDispatchApi 运营调度管理接口
type OperationDispatchApi struct{}

// CreateDispatch 创建运营调度
// @Tags     运营调度管理
// @Summary  创建运营调度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.OperationDispatch true "调度信息"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /operationDispatch/create [post]
func (a *OperationDispatchApi) CreateDispatch(c *gin.Context) {
	var req tripModel.OperationDispatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.OperationDispatchService.CreateDispatch(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateDispatch 更新运营调度
// @Tags     运营调度管理
// @Summary  更新运营调度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     tripModel.OperationDispatch true "调度信息"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /operationDispatch/update [put]
func (a *OperationDispatchApi) UpdateDispatch(c *gin.Context) {
	var req tripModel.OperationDispatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.OperationDispatchService.UpdateDispatch(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteDispatch 删除运营调度
// @Tags     运营调度管理
// @Summary  删除运营调度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.GetById true "主键ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /operationDispatch/delete [delete]
func (a *OperationDispatchApi) DeleteDispatch(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.OperationDispatchService.DeleteDispatch(req.Uint()); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDispatchByIds 批量删除运营调度
// @Tags     运营调度管理
// @Summary  批量删除运营调度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.IdsReq true "id列表"
// @Success  200  {object} response.Response{msg=string} "批量删除成功"
// @Router   /operationDispatch/deleteByIds [delete]
func (a *OperationDispatchApi) DeleteDispatchByIds(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.OperationDispatchService.DeleteDispatchByIds(req.Ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// CancelDispatch 撤销派单
// @Tags     运营调度管理
// @Summary  撤销派单（状态置为已取消）
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.GetById true "主键ID"
// @Success  200  {object} response.Response{msg=string} "撤销成功"
// @Router   /operationDispatch/cancel [post]
func (a *OperationDispatchApi) CancelDispatch(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.OperationDispatchService.CancelDispatch(req.Uint()); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("撤销派单成功", c)
}

// FindDispatch 根据id获取运营调度
// @Tags     运营调度管理
// @Summary  根据id获取运营调度
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    request.GetById true "主键ID"
// @Success  200  {object} response.Response{data=tripModel.OperationDispatch,msg=string} "获取成功"
// @Router   /operationDispatch/find [get]
func (a *OperationDispatchApi) FindDispatch(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := service.ServiceGroupApp.OperationDispatchService.GetDispatch(req.Uint())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetDispatchList 分页获取运营调度列表
// @Tags     运营调度管理
// @Summary  分页获取运营调度列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    tripReq.OperationDispatchSearch true "分页+筛选"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /operationDispatch/list [get]
func (a *OperationDispatchApi) GetDispatchList(c *gin.Context) {
	var info tripReq.OperationDispatchSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupApp.OperationDispatchService.GetDispatchList(info)
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
