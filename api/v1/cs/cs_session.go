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

type CsSessionApi struct{}

// CreateCsSession
// @Tags      CsSession
// @Summary   创建客服会话，会话编号由后端雪花ID自动生成，返回会话ID与会话编号供后续发消息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsSessionCreateReq                                    true  "关联工单, 用户ID, 用户姓名, 用户手机号, 客服ID, 客服姓名"
// @Success   200   {object}  response.Response{data=csRes.CsSessionCreateResponse,msg=string}  "创建客服会话"
// @Router    /api/cs/session [post]
func (a *CsSessionApi) CreateCsSession(c *gin.Context) {
	var req csReq.CsSessionCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsSessionCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	session, err := csSessionService.CreateCsSession(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("创建客服会话失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsSessionCreateResponse{
		SessionID: session.ID,
		SessionNo: session.SessionNo,
	}, "创建成功", c)
}

// DeleteCsSession
// @Tags      CsSession
// @Summary   删除客服会话（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "会话ID"
// @Success   200   {object}  response.Response{msg=string}  "删除客服会话"
// @Router    /api/cs/session [delete]
func (a *CsSessionApi) DeleteCsSession(c *gin.Context) {
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
	err = csSessionService.DeleteCsSession(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除客服会话失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsSession
// @Tags      CsSession
// @Summary   修改客服会话，会话编号不允许修改；状态：1进行中，2已结束
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsSessionUpdateReq          true  "会话ID, 会话信息, 状态, 开始时间, 结束时间"
// @Success   200   {object}  response.Response{msg=string}     "修改客服会话"
// @Router    /api/cs/session [put]
func (a *CsSessionApi) UpdateCsSession(c *gin.Context) {
	var req csReq.CsSessionUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsSessionUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csSessionService.UpdateCsSession(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("修改客服会话失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 会话状态等业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsSession
// @Tags      CsSession
// @Summary   获取客服会话详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                            true  "会话ID"
// @Success   200   {object}  response.Response{data=csRes.CsSessionResponse,msg=string}  "获取客服会话详情"
// @Router    /api/cs/session [get]
func (a *CsSessionApi) GetCsSession(c *gin.Context) {
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
	data, err := csSessionService.GetCsSession(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取客服会话失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsSessionResponse{Session: data}, "获取成功", c)
}

// GetCsSessionList
// @Tags      CsSession
// @Summary   分页获取客服会话列表，支持按状态、客服、用户与关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsSessionSearch                                     true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}     "分页获取客服会话列表"
// @Router    /api/cs/session/list [get]
func (a *CsSessionApi) GetCsSessionList(c *gin.Context) {
	var search csReq.CsSessionSearch
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
	list, total, err := csSessionService.GetCsSessionList(search)
	if err != nil {
		global.GVA_LOG.Error("获取客服会话列表失败!", zap.Error(err))
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
