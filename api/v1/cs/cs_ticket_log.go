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

type CsTicketLogApi struct{}

// CreateCsTicketLog
// @Tags      CsTicketLog
// @Summary   新增工单流水，操作人由后端从当前登录用户取，操作时间取记录创建时间
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsTicketLogCreateReq                                  true  "关联工单ID, 变更前状态, 变更后状态, 操作备注"
// @Success   200   {object}  response.Response{data=csRes.CsTicketLogResponse,msg=string}  "新增工单流水"
// @Router    /api/cs/ticketLog [post]
func (a *CsTicketLogApi) CreateCsTicketLog(c *gin.Context) {
	var req csReq.CsTicketLogCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsTicketLogCreateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	operatorID := utils.GetUserID(c)
	data, err := csTicketLogService.CreateCsTicketLog(operatorID, req)
	if err != nil {
		global.GVA_LOG.Error("新增工单流水失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(csRes.CsTicketLogResponse{Log: csRes.ToCsTicketLogDto(data)}, "创建成功", c)
}

// GetCsTicketLog
// @Tags      CsTicketLog
// @Summary   获取工单流水详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                            true  "工单流水ID"
// @Success   200   {object}  response.Response{data=csRes.CsTicketLogResponse,msg=string}  "获取工单流水详情"
// @Router    /api/cs/ticketLog [get]
func (a *CsTicketLogApi) GetCsTicketLog(c *gin.Context) {
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
	data, err := csTicketLogService.GetCsTicketLog(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取工单流水失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsTicketLogResponse{Log: csRes.ToCsTicketLogDto(data)}, "获取成功", c)
}

// GetCsTicketLogList
// @Tags      CsTicketLog
// @Summary   分页获取工单流水列表，支持按工单ID、操作人、变更后状态与关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsTicketLogSearch                                   true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}     "分页获取工单流水列表"
// @Router    /api/cs/ticketLog/list [get]
func (a *CsTicketLogApi) GetCsTicketLogList(c *gin.Context) {
	var search csReq.CsTicketLogSearch
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
	list, total, err := csTicketLogService.GetCsTicketLogList(search)
	if err != nil {
		global.GVA_LOG.Error("获取工单流水列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     csRes.ToCsTicketLogDtoList(list),
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}
