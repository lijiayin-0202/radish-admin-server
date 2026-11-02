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

type CsNoticeApi struct{}

// CreateCsNotice
// @Tags      CsNotice
// @Summary   新增站内公告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNoticeCreateReq                                    true  "公告标题, 公告内容, 公告类型, 状态, 排序"
// @Success   200   {object}  response.Response{data=csRes.CsNoticeResponse,msg=string}   "新增站内公告"
// @Router    /api/cs/notice [post]
func (a *CsNoticeApi) CreateCsNotice(c *gin.Context) {
	var req csReq.CsNoticeCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsNoticeCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	data, err := csNoticeService.CreateCsNotice(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("新增站内公告失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsNoticeResponse{Notice: csRes.ToCsNoticeDto(data)}, "创建成功", c)
}

// DeleteCsNotice
// @Tags      CsNotice
// @Summary   删除站内公告（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "站内公告ID"
// @Success   200   {object}  response.Response{msg=string}  "删除站内公告"
// @Router    /api/cs/notice [delete]
func (a *CsNoticeApi) DeleteCsNotice(c *gin.Context) {
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
	err = csNoticeService.DeleteCsNotice(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除站内公告失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsNotice
// @Tags      CsNotice
// @Summary   修改站内公告
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNoticeUpdateReq        true  "站内公告ID, 公告标题, 公告内容, 公告类型, 状态, 排序"
// @Success   200   {object}  response.Response{msg=string}  "修改站内公告"
// @Router    /api/cs/notice [put]
func (a *CsNoticeApi) UpdateCsNotice(c *gin.Context) {
	var req csReq.CsNoticeUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsNoticeUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	err = csNoticeService.UpdateCsNotice(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("修改站内公告失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsNotice
// @Tags      CsNotice
// @Summary   获取站内公告详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                          true  "站内公告ID"
// @Success   200   {object}  response.Response{data=csRes.CsNoticeResponse,msg=string}  "获取站内公告详情"
// @Router    /api/cs/notice [get]
func (a *CsNoticeApi) GetCsNotice(c *gin.Context) {
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
	data, err := csNoticeService.GetCsNotice(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取站内公告失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsNoticeResponse{Notice: csRes.ToCsNoticeDto(data)}, "获取成功", c)
}

// GetCsNoticeList
// @Tags      CsNotice
// @Summary   分页获取站内公告列表，支持按标题、关键词、类型与状态筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsNoticeSearch                                 true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取站内公告列表"
// @Router    /api/cs/notice/list [get]
func (a *CsNoticeApi) GetCsNoticeList(c *gin.Context) {
	var search csReq.CsNoticeSearch
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
	list, total, err := csNoticeService.GetCsNoticeList(search)
	if err != nil {
		global.GVA_LOG.Error("获取站内公告列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     csRes.ToCsNoticeDtoList(list),
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}
