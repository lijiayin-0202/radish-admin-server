package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// CsSearch 客服业务通用查询
type CsSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// TicketSearch 工单查询
type TicketSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// CreateTicketReq 创建工单
type CreateTicketReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  *int   `json:"status"`
}

// UpdateTicketStatusReq 修改工单状态
type UpdateTicketStatusReq struct {
	ID     int `json:"id" binding:"required"`
	Status int `json:"status" binding:"required"`
}
