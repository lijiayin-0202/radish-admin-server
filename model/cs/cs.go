package cs

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ==================== 客服中心（模块5，luobo 库） ====================

// CsSession 客服会话
type CsSession struct {
	global.GVA_MODEL
	SessionNo     string     `json:"sessionNo" gorm:"type:varchar(50);comment:会话编号"`
	PassengerID   uint       `json:"passengerId" gorm:"comment:乘客ID"`
	OrderID       uint       `json:"orderId" gorm:"comment:关联订单ID"`
	Channel       string     `json:"channel" gorm:"type:varchar(50);comment:渠道 在线/电话"`
	Scene         string     `json:"scene" gorm:"type:varchar(50);comment:场景"`
	IssueCategory string     `json:"issueCategory" gorm:"type:varchar(100);comment:问题分类"`
	SessionStatus string     `json:"sessionStatus" gorm:"type:varchar(50);comment:状态 排队中/服务中/已结束"`
	AssignedTo    string     `json:"assignedTo" gorm:"type:varchar(100);comment:接入客服"`
	ConnectedAt   *time.Time `json:"connectedAt" gorm:"comment:接入时间"`
	ClosedAt      *time.Time `json:"closedAt" gorm:"comment:结束时间"`
}

func (CsSession) TableName() string { return "cs_sessions" }

// CsTicket 客服工单
type CsTicket struct {
	global.GVA_MODEL
	TicketNo    string     `json:"ticketNo" gorm:"type:varchar(50);comment:工单号"`
	PassengerID uint       `json:"passengerId" gorm:"comment:乘客ID"`
	OrderID     uint       `json:"orderId" gorm:"comment:关联订单ID"`
	TicketType  string     `json:"ticketType" gorm:"type:varchar(50);comment:类型 咨询/投诉/退款/故障反馈/合作咨询"`
	Priority    string     `json:"priority" gorm:"type:varchar(20);comment:优先级"`
	Title       string     `json:"title" gorm:"type:varchar(200);comment:标题"`
	Description string     `json:"description" gorm:"type:text;comment:描述"`
	Status      int        `json:"status" gorm:"type:int;default:1;comment:状态 1待处理/2处理中/3已关闭"`
	AssignedTo  string     `json:"assignedTo" gorm:"type:varchar(100);comment:处理人"`
	Resolution  string     `json:"resolution" gorm:"type:varchar(300);comment:处理结果"`
	ClosedAt    *time.Time `json:"closedAt" gorm:"comment:关闭时间"`
}

func (CsTicket) TableName() string { return "cs_tickets" }

// CsFaqCategory FAQ 分类
type CsFaqCategory struct {
	global.GVA_MODEL
	CategoryName string `json:"categoryName" gorm:"type:varchar(100);comment:分类名称"`
	Status       int    `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
	Sort         int    `json:"sort" gorm:"type:int;default:0;comment:排序"`
}

func (CsFaqCategory) TableName() string { return "cs_faq_categories" }

// CsFaq FAQ 问答
type CsFaq struct {
	global.GVA_MODEL
	CategoryID uint   `json:"categoryId" gorm:"comment:分类ID"`
	Question   string `json:"question" gorm:"type:varchar(300);comment:问题"`
	Answer     string `json:"answer" gorm:"type:text;comment:答案"`
	Keywords   string `json:"keywords" gorm:"type:varchar(300);comment:检索关键词"`
	Status     int    `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
	Sort       int    `json:"sort" gorm:"type:int;default:0;comment:排序"`
}

func (CsFaq) TableName() string { return "cs_faqs" }

// CsNotification 消息通知
type CsNotification struct {
	global.GVA_MODEL
	PassengerID uint   `json:"passengerId" gorm:"comment:乘客ID"`
	NotifyType  string `json:"notifyType" gorm:"type:varchar(50);comment:类型 订单/支付/活动/安全/系统"`
	Title       string `json:"title" gorm:"type:varchar(200);comment:标题"`
	Content     string `json:"content" gorm:"type:text;comment:内容"`
	IsRead      bool   `json:"isRead" gorm:"comment:是否已读"`
}

func (CsNotification) TableName() string { return "cs_notifications" }

// CsTicketLog 工单处理日志（模块5，PRD 的 service_ticket_logs，本次补建）
type CsTicketLog struct {
	global.GVA_MODEL
	TicketID uint   `json:"ticketId" gorm:"comment:工单ID"`
	Action   string `json:"action" gorm:"type:varchar(50);comment:动作 创建/回复/转派/关闭"`
	Operator string `json:"operator" gorm:"type:varchar(100);comment:操作人"`
	Content  string `json:"content" gorm:"type:varchar(500);comment:操作内容"`
}

func (CsTicketLog) TableName() string { return "cs_ticket_logs" }
