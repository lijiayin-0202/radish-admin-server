package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	platformModel "github.com/flipped-aurora/gin-vue-admin/server/model/platform"
	platformReq "github.com/flipped-aurora/gin-vue-admin/server/model/platform/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
)

// PricingRuleApi 计价规则接口
type PricingRuleApi struct{}

// CreatePricingRule 创建计价规则
// @Tags     计价规则
// @Summary  创建计价规则
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     platformModel.PricingRule true "计价规则"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /pricingRule/create [post]
func (a *PricingRuleApi) CreatePricingRule(c *gin.Context) {
	var req platformModel.PricingRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.PricingRuleService.CreatePricingRule(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdatePricingRule 更新计价规则
// @Tags     计价规则
// @Summary  更新计价规则
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     platformModel.PricingRule true "计价规则"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /pricingRule/update [put]
func (a *PricingRuleApi) UpdatePricingRule(c *gin.Context) {
	var req platformModel.PricingRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.PricingRuleService.UpdatePricingRule(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeletePricingRule 删除计价规则
// @Tags     计价规则
// @Summary  删除计价规则
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.GetById true "主键ID"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /pricingRule/delete [delete]
func (a *PricingRuleApi) DeletePricingRule(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.PricingRuleService.DeletePricingRule(req.Uint()); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeletePricingRulesByIds 批量删除计价规则
// @Tags     计价规则
// @Summary  批量删除计价规则
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body     request.IdsReq true "id列表"
// @Success  200  {object} response.Response{msg=string} "批量删除成功"
// @Router   /pricingRule/deleteByIds [delete]
func (a *PricingRuleApi) DeletePricingRulesByIds(c *gin.Context) {
	var req request.IdsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.ServiceGroupApp.PricingRuleService.DeletePricingRulesByIds(req.Ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// FindPricingRule 根据id获取计价规则
// @Tags     计价规则
// @Summary  根据id获取计价规则
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    request.GetById true "主键ID"
// @Success  200  {object} response.Response{data=platformModel.PricingRule,msg=string} "获取成功"
// @Router   /pricingRule/find [get]
func (a *PricingRuleApi) FindPricingRule(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := service.ServiceGroupApp.PricingRuleService.GetPricingRule(req.Uint())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetPricingRuleList 分页获取计价规则列表
// @Tags     计价规则
// @Summary  分页获取计价规则列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query    platformReq.PricingRuleSearch true "分页+筛选"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /pricingRule/list [get]
func (a *PricingRuleApi) GetPricingRuleList(c *gin.Context) {
	var info platformReq.PricingRuleSearch
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupApp.PricingRuleService.GetPricingRuleList(info)
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
