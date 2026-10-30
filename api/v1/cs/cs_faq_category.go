package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	csRes "github.com/flipped-aurora/gin-vue-admin/server/model/cs/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CsFaqCategoryApi struct{}

// CreateCsFaqCategory
// @Tags      CsFaqCategory
// @Summary   新增常见问题分类，status 传 0 禁用 / 1 启用，不传按默认1启用
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsFaqCategoryCreateReq                                    true  "分类名称, 排序号, 状态, 备注"
// @Success   200   {object}  response.Response{data=csRes.CsFaqCategoryResponse,msg=string}   "新增常见问题分类"
// @Router    /api/cs/faqCategory [post]
func (a *CsFaqCategoryApi) CreateCsFaqCategory(c *gin.Context) {
	var req csReq.CsFaqCategoryCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsFaqCategoryCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	data, err := csFaqCategoryService.CreateCsFaqCategory(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("新增常见问题分类失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsFaqCategoryResponse{Category: data}, "创建成功", c)
}

// DeleteCsFaqCategory
// @Tags      CsFaqCategory
// @Summary   删除常见问题分类（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "分类ID"
// @Success   200   {object}  response.Response{msg=string}  "删除常见问题分类"
// @Router    /api/cs/faqCategory [delete]
func (a *CsFaqCategoryApi) DeleteCsFaqCategory(c *gin.Context) {
	var idReq request.GetById
	err := c.ShouldBindJSON(&idReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idReq, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csFaqCategoryService.DeleteCsFaqCategory(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除常见问题分类失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsFaqCategory
// @Tags      CsFaqCategory
// @Summary   修改常见问题分类，支持启停：status 传 0 禁用 / 1 启用，不传表示不变更
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsFaqCategoryUpdateReq   true  "分类ID, 分类名称, 排序号, 状态, 备注"
// @Success   200   {object}  response.Response{msg=string}  "修改常见问题分类"
// @Router    /api/cs/faqCategory [put]
func (a *CsFaqCategoryApi) UpdateCsFaqCategory(c *gin.Context) {
	var req csReq.CsFaqCategoryUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsFaqCategoryUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	err = csFaqCategoryService.UpdateCsFaqCategory(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("修改常见问题分类失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 启用状态等业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsFaqCategory
// @Tags      CsFaqCategory
// @Summary   获取常见问题分类详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                                true  "分类ID"
// @Success   200   {object}  response.Response{data=csRes.CsFaqCategoryResponse,msg=string}  "获取常见问题分类详情"
// @Router    /api/cs/faqCategory [get]
func (a *CsFaqCategoryApi) GetCsFaqCategory(c *gin.Context) {
	var idReq request.GetById
	err := c.ShouldBindQuery(&idReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idReq, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := csFaqCategoryService.GetCsFaqCategory(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取常见问题分类失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsFaqCategoryResponse{Category: data}, "获取成功", c)
}

// GetCsFaqCategoryList
// @Tags      CsFaqCategory
// @Summary   分页获取常见问题分类列表，支持按分类名称与启用状态筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsFaqCategorySearch                                 true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}     "分页获取常见问题分类列表"
// @Router    /api/cs/faqCategory/list [get]
func (a *CsFaqCategoryApi) GetCsFaqCategoryList(c *gin.Context) {
	var search csReq.CsFaqCategorySearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(search.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := csFaqCategoryService.GetCsFaqCategoryList(search)
	if err != nil {
		global.GVA_LOG.Error("获取常见问题分类列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}
