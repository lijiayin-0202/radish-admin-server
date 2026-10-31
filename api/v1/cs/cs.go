package cs

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	csService "github.com/flipped-aurora/gin-vue-admin/server/service/cs"
	"github.com/gin-gonic/gin"
)

type CsApi struct{}

// ==================== 会话 ====================

func (a *CsApi) GetSessionList(c *gin.Context) {
	var info csReq.CsSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := csService.CsServiceApp.GetSessionList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// ==================== 工单 ====================

func (a *CsApi) GetTicketList(c *gin.Context) {
	var info csReq.TicketSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := csService.CsServiceApp.GetTicketList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *CsApi) UpdateTicket(c *gin.Context) {
	var req csModel.CsTicket
	_ = c.ShouldBindJSON(&req)
	if err := csService.CsServiceApp.UpdateTicket(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ==================== 工单（新增接口） ====================

func (a *CsApi) CreateTicket(c *gin.Context) {
	var req csReq.CreateTicketReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	st := 1
	if req.Status != nil {
		st = *req.Status
	}
	ticket := csModel.CsTicket{
		Title:       req.Title,
		Description: req.Content,
		Status:      st,
	}
	if err := csService.CsServiceApp.CreateTicket(&ticket); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"id": ticket.ID}, "创建成功", c)
}

func (a *CsApi) GetTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("id 参数非法", c)
		return
	}
	ticket, err := csService.CsServiceApp.GetTicket(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(ticket, "获取成功", c)
}

func (a *CsApi) UpdateTicketStatus(c *gin.Context) {
	var req csReq.UpdateTicketStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := csService.CsServiceApp.UpdateTicketStatus(req.ID, req.Status); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ==================== FAQ 分类 ====================

func (a *CsApi) GetFaqCategoryList(c *gin.Context) {
	var info csReq.CsSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := csService.CsServiceApp.GetFaqCategoryList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *CsApi) CreateFaqCategory(c *gin.Context) {
	var req csModel.CsFaqCategory
	_ = c.ShouldBindJSON(&req)
	if err := csService.CsServiceApp.CreateFaqCategory(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *CsApi) UpdateFaqCategory(c *gin.Context) {
	var req csModel.CsFaqCategory
	_ = c.ShouldBindJSON(&req)
	if err := csService.CsServiceApp.UpdateFaqCategory(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *CsApi) DeleteFaqCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := csService.CsServiceApp.DeleteFaqCategory(uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== FAQ ====================

func (a *CsApi) GetFaqList(c *gin.Context) {
	var info csReq.CsSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := csService.CsServiceApp.GetFaqList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *CsApi) CreateFaq(c *gin.Context) {
	var req csModel.CsFaq
	_ = c.ShouldBindJSON(&req)
	if err := csService.CsServiceApp.CreateFaq(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *CsApi) UpdateFaq(c *gin.Context) {
	var req csModel.CsFaq
	_ = c.ShouldBindJSON(&req)
	if err := csService.CsServiceApp.UpdateFaq(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *CsApi) DeleteFaq(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := csService.CsServiceApp.DeleteFaq(uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== 消息通知 ====================

func (a *CsApi) GetNotificationList(c *gin.Context) {
	var info csReq.CsSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := csService.CsServiceApp.GetNotificationList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}
