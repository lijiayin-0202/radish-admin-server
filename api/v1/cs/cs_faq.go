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

type CsFaqApi struct{}

// CreateCsFaq
// @Tags      CsFaq
// @Summary   新增常见问题，status 传 0 禁用 / 1 启用，不传按默认1启用
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsFaqCreateReq                                  true  "问题标题, 回答内容, 问题分类, 标签, 排序号, 状态"
// @Success   200   {object}  response.Response{data=csRes.CsFaqResponse,msg=string}  "新增常见问题"
// @Router    /api/cs/faq [post]
func (a *CsFaqApi) CreateCsFaq(c *gin.Context) {
	var req csReq.CsFaqCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsFaqCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := csFaqService.CreateCsFaq(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("新增常见问题失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsFaqResponse{Faq: data}, "创建成功", c)
}

// DeleteCsFaq
// @Tags      CsFaq
// @Summary   删除常见问题（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "FAQID"
// @Success   200   {object}  response.Response{msg=string}  "删除常见问题"
// @Router    /api/cs/faq [delete]
func (a *CsFaqApi) DeleteCsFaq(c *gin.Context) {
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
	err = csFaqService.DeleteCsFaq(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除常见问题失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsFaq
// @Tags      CsFaq
// @Summary   修改常见问题，支持上下架：status 传 0 下架 / 1 上架，不传表示不变更
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsFaqUpdateReq            true  "FAQID, 问题标题, 回答内容, 问题分类, 标签, 排序号, 上架状态"
// @Success   200   {object}  response.Response{msg=string}   "修改常见问题"
// @Router    /api/cs/faq [put]
func (a *CsFaqApi) UpdateCsFaq(c *gin.Context) {
	var req csReq.CsFaqUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsFaqUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csFaqService.UpdateCsFaq(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("修改常见问题失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 上下架状态等业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsFaq
// @Tags      CsFaq
// @Summary   获取常见问题详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                         true  "FAQID"
// @Success   200   {object}  response.Response{data=csRes.CsFaqResponse,msg=string}  "获取常见问题详情"
// @Router    /api/cs/faq [get]
func (a *CsFaqApi) GetCsFaq(c *gin.Context) {
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
	data, err := csFaqService.GetCsFaq(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取常见问题失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsFaqResponse{Faq: data}, "获取成功", c)
}

// GetCsFaqList
// @Tags      CsFaq
// @Summary   分页获取常见问题列表，支持按分类、上架状态与关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsFaqSearch                                      true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取常见问题列表"
// @Router    /api/cs/faq/list [get]
func (a *CsFaqApi) GetCsFaqList(c *gin.Context) {
	var search csReq.CsFaqSearch
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
	list, total, err := csFaqService.GetCsFaqList(search)
	if err != nil {
		global.GVA_LOG.Error("获取常见问题列表失败!", zap.Error(err))
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
