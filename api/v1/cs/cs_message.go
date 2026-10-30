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

type CsMessageApi struct{}

// CreateCsMessage
// @Tags      CsMessage
// @Summary   保存客服会话消息，消息序号与发送时间不传时由后端自动生成
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsMessageCreateReq                              true  "归属会话ID, 发送方类型, 消息类型, 消息文本内容, 附件地址"
// @Success   200   {object}  response.Response{data=csRes.CsMessageResponse,msg=string}  "保存客服会话消息"
// @Router    /api/cs/message [post]
func (a *CsMessageApi) CreateCsMessage(c *gin.Context) {
	var req csReq.CsMessageCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsMessageCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	data, err := csMessageService.CreateCsMessage(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("保存客服会话消息失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsMessageResponse{Message: data}, "创建成功", c)
}

// DeleteCsMessage
// @Tags      CsMessage
// @Summary   删除客服会话消息（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "消息ID"
// @Success   200   {object}  response.Response{msg=string}  "删除客服会话消息"
// @Router    /api/cs/message [delete]
func (a *CsMessageApi) DeleteCsMessage(c *gin.Context) {
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
	err = csMessageService.DeleteCsMessage(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除客服会话消息失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsMessage
// @Tags      CsMessage
// @Summary   修改客服会话消息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsMessageUpdateReq         true  "消息ID, 发送方类型, 消息类型, 消息文本内容, 附件地址, 消息序号, 已读状态"
// @Success   200   {object}  response.Response{msg=string}    "修改客服会话消息"
// @Router    /api/cs/message [put]
func (a *CsMessageApi) UpdateCsMessage(c *gin.Context) {
	var req csReq.CsMessageUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsMessageUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	err = csMessageService.UpdateCsMessage(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("修改客服会话消息失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsMessage
// @Tags      CsMessage
// @Summary   获取客服会话消息详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                            true  "消息ID"
// @Success   200   {object}  response.Response{data=csRes.CsMessageResponse,msg=string}  "获取客服会话消息详情"
// @Router    /api/cs/message [get]
func (a *CsMessageApi) GetCsMessage(c *gin.Context) {
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
	data, err := csMessageService.GetCsMessage(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取客服会话消息失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsMessageResponse{Message: data}, "获取成功", c)
}

// GetCsMessageList
// @Tags      CsMessage
// @Summary   分页获取客服会话消息列表，传 sessionId 时按会话历史消息语义返回（消息序号正序）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsMessageSearch                                   true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}   "分页获取客服会话消息列表"
// @Router    /api/cs/message/list [get]
func (a *CsMessageApi) GetCsMessageList(c *gin.Context) {
	var search csReq.CsMessageSearch
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
	list, total, err := csMessageService.GetCsMessageList(search)
	if err != nil {
		global.GVA_LOG.Error("获取客服会话消息列表失败!", zap.Error(err))
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
