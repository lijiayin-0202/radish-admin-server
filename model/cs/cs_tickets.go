package cs

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsTicket 客服工单表
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
// TODO(ride-hailing)：order_id / order_no 仅做关联存储，订单与车辆信息由 ride-hailing 微服务查询接口提供，
// 本仓库不落地订单/车辆 model
type CsTicket struct {
	global.GVA_MODEL
	TicketNo          string     `json:"ticketNo" form:"ticketNo" gorm:"column:ticket_no;type:varchar(64);not null;default:'';comment:工单号"`
	UserID            int64      `json:"userId" form:"userId" gorm:"column:user_id;type:bigint;not null;default:0;comment:用户ID"`
	UserName          string     `json:"userName" form:"userName" gorm:"column:user_name;type:varchar(128);not null;default:'';comment:用户姓名"`
	UserPhone         string     `json:"userPhone" form:"userPhone" gorm:"column:user_phone;type:varchar(32);not null;default:'';comment:用户手机号"`
	OrderID           int64      `json:"orderId" form:"orderId" gorm:"column:order_id;type:bigint;not null;default:0;comment:关联订单ID"`
	OrderNo           string     `json:"orderNo" form:"orderNo" gorm:"column:order_no;type:varchar(64);not null;default:'';comment:关联订单号"`
	Status            int64      `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:工单状态：1待处理，2处理中，3已解决，4已关闭"`
	Priority          int64      `json:"priority" form:"priority" gorm:"column:priority;not null;default:1;comment:优先级：1低，2中，3高"`
	Category          string     `json:"category" form:"category" gorm:"column:category;type:varchar(128);not null;default:'';comment:问题分类"`
	Title             string     `json:"title" form:"title" gorm:"column:title;type:varchar(256);not null;default:'';comment:工单标题"`
	Description       string     `json:"description" form:"description" gorm:"column:description;type:text;comment:问题描述"`
	AttachUrls        string     `json:"attachUrls" form:"attachUrls" gorm:"column:attach_urls;type:text;comment:附件URL列表"`
	AssignedClerkID   int64      `json:"assignedClerkId" form:"assignedClerkId" gorm:"column:assigned_clerk_id;type:bigint;not null;default:0;comment:分配客服ID"`
	AssignedClerkName string     `json:"assignedClerkName" form:"assignedClerkName" gorm:"column:assigned_clerk_name;type:varchar(128);not null;default:'';comment:分配客服姓名"`
	ResolvedAt        *time.Time `json:"resolvedAt" form:"resolvedAt" gorm:"column:resolved_at;comment:解决时间"`
	ClosedAt          *time.Time `json:"closedAt" form:"closedAt" gorm:"column:closed_at;comment:关闭时间"`
	Remark            string     `json:"remark" form:"remark" gorm:"column:remark;type:text;comment:备注"`
	CreatedBy         int64      `json:"createdBy" form:"createdBy" gorm:"column:created_by;type:bigint;not null;default:0;comment:创建者"`
}

// TableName CsTicket 自定义表名 cs_tickets
func (CsTicket) TableName() string {
	return "cs_tickets"
}
