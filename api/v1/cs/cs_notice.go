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
// @Summary   新增消息通知，新增固定为未读，已读状态由 /read 接口标记
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNoticeCreateReq                                    true  "接收人, 消息标题, 消息内容, 消息类型"
// @Success   200   {object}  response.Response{data=csRes.CsNoticeResponse,msg=string}   "新增消息通知"
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
		global.GVA_LOG.Error("新增消息通知失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsNoticeResponse{Notice: data}, "创建成功", c)
}

// DeleteCsNotice
// @Tags      CsNotice
// @Summary   删除消息通知（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "消息通知ID"
// @Success   200   {object}  response.Response{msg=string}  "删除消息通知"
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
		global.GVA_LOG.Error("删除消息通知失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsNotice
// @Tags      CsNotice
// @Summary   修改消息通知，已读状态建议走 /read 接口标记
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNoticeUpdateReq        true  "消息通知ID, 接收人, 消息标题, 消息内容, 消息类型, 已读状态"
// @Success   200   {object}  response.Response{msg=string}  "修改消息通知"
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
		global.GVA_LOG.Error("修改消息通知失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 消息类型等业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ReadCsNotice
// @Tags      CsNotice
// @Summary   标记消息通知已读，重复标记幂等返回
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNoticeReadReq          true  "消息通知ID"
// @Success   200   {object}  response.Response{msg=string}  "标记消息通知已读"
// @Router    /api/cs/notice/read [put]
func (a *CsNoticeApi) ReadCsNotice(c *gin.Context) {
	var req csReq.CsNoticeReadReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsNoticeReadVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csNoticeService.ReadCsNotice(req.ID)
	if err != nil {
		global.GVA_LOG.Error("标记消息通知已读失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsNotice
// @Tags      CsNotice
// @Summary   获取消息通知详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                          true  "消息通知ID"
// @Success   200   {object}  response.Response{data=csRes.CsNoticeResponse,msg=string}  "获取消息通知详情"
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
		global.GVA_LOG.Error("获取消息通知失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsNoticeResponse{Notice: data}, "获取成功", c)
}

// GetCsNoticeList
// @Tags      CsNotice
// @Summary   分页获取消息通知列表，支持按接收人、已读状态、消息类型与关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsNoticeSearch                                 true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取消息通知列表"
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
		global.GVA_LOG.Error("获取消息通知列表失败!", zap.Error(err))
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
