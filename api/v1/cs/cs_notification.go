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

type CsNotificationApi struct{}

// GetCsNotificationList
// @Tags      CsNotification
// @Summary   分页获取我的通知列表，仅返回当前登录用户的通知，支持已读状态、类型与关键词筛选
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     csReq.CsNotificationSearch                              true  "页码, 每页大小, 查询条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取通知列表"
// @Router    /api/cs/notification/list [get]
func (a *CsNotificationApi) GetCsNotificationList(c *gin.Context) {
	var search csReq.CsNotificationSearch
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
	// 业务规则：通知列表仅返回当前登录用户自己的通知，禁止越权查看他人通知
	operatorID := utils.GetUserID(c)
	uid := int64(operatorID)
	search.UserID = &uid
	list, total, err := csNotificationService.GetCsNotificationList(search)
	if err != nil {
		global.GVA_LOG.Error("获取通知列表失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     csRes.ToCsNotificationDtoList(list),
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// ReadCsNotification
// @Tags      CsNotification
// @Summary   标记通知已读，重复标记幂等返回
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      csReq.CsNotificationReadReq   true  "通知ID"
// @Success   200   {object}  response.Response{msg=string}  "标记通知已读"
// @Router    /api/cs/notification/read [put]
func (a *CsNotificationApi) ReadCsNotification(c *gin.Context) {
	var req csReq.CsNotificationReadReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, CsNotificationReadVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = csNotificationService.ReadCsNotification(req.ID)
	if err != nil {
		global.GVA_LOG.Error("标记通知已读失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetCsNotification
// @Tags      CsNotification
// @Summary   获取通知详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.GetById                                            true  "通知ID"
// @Success   200   {object}  response.Response{data=csRes.CsNotificationResponse,msg=string}  "获取通知详情"
// @Router    /api/cs/notification [get]
func (a *CsNotificationApi) GetCsNotification(c *gin.Context) {
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
	data, err := csNotificationService.GetCsNotification(idReq.Uint())
	if err != nil {
		global.GVA_LOG.Error("获取通知失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(csRes.CsNotificationResponse{Notification: csRes.ToCsNotificationDto(data)}, "获取成功", c)
}
