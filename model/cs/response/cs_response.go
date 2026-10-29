package response

import (
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
)

// CsTicketResponse 客服工单详情返回
type CsTicketResponse struct {
	Ticket csModel.CsTicket `json:"ticket"`
}

// CsSessionResponse 客服会话详情返回
type CsSessionResponse struct {
	Session csModel.CsSession `json:"session"`
}

// CsSessionCreateResponse 创建客服会话返回，会话编号由后端生成，前端据此发起后续会话消息
type CsSessionCreateResponse struct {
	SessionID uint   `json:"sessionId"`
	SessionNo string `json:"sessionNo"`
}

// CsMessageResponse 客服会话消息详情返回
type CsMessageResponse struct {
	Message csModel.CsMessage `json:"message"`
}

// CsFaqResponse 常见问题详情返回
type CsFaqResponse struct {
	Faq csModel.CsFaq `json:"faq"`
}
