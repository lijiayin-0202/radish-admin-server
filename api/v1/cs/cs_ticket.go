package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	csRes "github.com/flipped-aurora/gin-vue-admin/server/model/cs/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CsTicketApi struct{}

// CreateCsTicket
// @Tags      CsTicket
// @Summary   新增客服工单，工单号由后端雪花ID自动生成，新建固定为待受理
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsTicketCreateReq                              true  "用户ID, 用户姓名, 用户手机号, 关联订单, 优先级, 问题分类, 工单标题, 问题描述"
// @Success   200   {object}  response.Response{data=csRes.CsTicketResponse,msg=string}  "新增客服工单"
// @Router    /api/cs/ticket [post]
func (a *CsTicketApi) CreateCsTicket(c *gin.Context) {
	var req csReq.CsTicketCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsTicketCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := csTicketService.CreateCsTicket(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("新增客服工单失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsTicketResponse{Ticket: data}, "创建成功", c)
}

// DeleteCsTicket
// @Tags      CsTicket
// @Summary   删除客服工单（软删除）
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      uint                           true  "工单ID"
// @Success   200   {object}  response.Response{msg=string}  "删除客服工单"
// @Router    /api/cs/ticket/{id} [delete]
func (a *CsTicketApi) DeleteCsTicket(c *gin.Context) {
	var idReq csReq.CsTicketIdReq
	err := c.ShouldBindUri(&idReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idReq, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csTicketService.DeleteCsTicket(utils.GetUserID(c), idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("删除客服工单失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCsTicket
// @Tags      CsTicket
// @Summary   修改客服工单，工单号不允许修改；状态流转：待受理(1) -> 处理中(2) -> 已完结(3)
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsTicketUpdateReq          true  "工单ID, 工单信息, 状态, 分配客服, 解决时间, 关闭时间"
// @Success   200   {object}  response.Response{msg=string}     "修改客服工单"
// @Router    /api/cs/ticket [put]
func (a *CsTicketApi) UpdateCsTicket(c *gin.Context) {
	var req csReq.CsTicketUpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsTicketUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csTicketService.UpdateCsTicket(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("修改客服工单失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 状态流转等业务校验错误直接回传提示
		return
	}
	response.OkWithMessage("更新成功", c)
}

// UpdateCsTicketStatus
// @Tags      CsTicket
// @Summary   修改客服工单状态，按流转规则校验：待受理(1) -> 处理中(2) -> 已完结(3)
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsTicketStatusReq                                    true  "工单ID, 目标状态, 分配客服, 解决时间, 关闭时间, 备注"
// @Success   200   {object}  response.Response{data=csRes.CsTicketResponse,msg=string}  "修改客服工单状态"
// @Router    /api/cs/ticket/status [put]
func (a *CsTicketApi) UpdateCsTicketStatus(c *gin.Context) {
	var req csReq.CsTicketStatusReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsTicketStatusVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	entity, err := csTicketService.UpdateCsTicketStatus(utils.GetUserID(c), req)
	if err != nil {
		global.GVA_LOG.Error("修改客服工单状态失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c) // 状态流转等业务校验错误直接回传提示
		return
	}
	response.OkWithDetailed(csRes.CsTicketResponse{Ticket: entity}, "更新成功", c)
}

// GetCsTicket
// @Tags      CsTicket
// @Summary   获取客服工单详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id   path      uint                                                        true  "工单ID"
// @Success   200   {object}  response.Response{data=csRes.CsTicketResponse,msg=string}  "获取客服工单详情"
// @Router    /api/cs/ticket/{id} [get]
func (a *CsTicketApi) GetCsTicket(c *gin.Context) {
	var idReq csReq.CsTicketIdReq
	err := c.ShouldBindUri(&idReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idReq, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := csTicketService.GetCsTicket(utils.GetUserID(c), idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取客服工单失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsTicketResponse{Ticket: data}, "获取成功", c)
}

// GetCsTicketList
// @Tags      CsTicket
// @Summary   分页获取客服工单列表，支持按状态、优先级、工单号、关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsTicketSearch                                     true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}    "分页获取客服工单列表"
// @Router    /api/cs/ticket/list [get]
func (a *CsTicketApi) GetCsTicketList(c *gin.Context) {
	var search csReq.CsTicketSearch
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
	list, total, err := csTicketService.GetCsTicketList(search)
	if err != nil {
		global.GVA_LOG.Error("获取客服工单列表失败!", zap.Error(err))
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
