package cs

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsSession 客服会话表
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
// TODO(ride-hailing)：会话如需展示订单/车辆信息，调用 ride-hailing 微服务查询接口获取，本仓库不落地订单/车辆 model
type CsSession struct {
	global.GVA_MODEL
	SessionNo      string     `json:"sessionNo" form:"sessionNo" gorm:"column:session_no;type:varchar(64);not null;default:'';comment:会话编号"`
	TicketID       int64      `json:"ticketId" form:"ticketId" gorm:"column:ticket_id;type:bigint;not null;default:0;comment:关联工单ID"`
	UserID         int64      `json:"userId" form:"userId" gorm:"column:user_id;type:bigint;not null;default:0;comment:用户ID"`
	UserName       string     `json:"userName" form:"userName" gorm:"column:user_name;type:varchar(128);not null;default:'';comment:用户姓名"`
	UserPhone      string     `json:"userPhone" form:"userPhone" gorm:"column:user_phone;type:varchar(32);not null;default:'';comment:用户手机号"`
	ClerkID        int64      `json:"clerkId" form:"clerkId" gorm:"column:clerk_id;type:bigint;not null;default:0;comment:客服ID"`
	ClerkName      string     `json:"clerkName" form:"clerkName" gorm:"column:clerk_name;type:varchar(128);not null;default:'';comment:客服姓名"`
	Status         int64      `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:会话状态：1进行中，2已结束"`
	StartedAt      *time.Time `json:"startedAt" form:"startedAt" gorm:"column:started_at;comment:开始时间"`
	EndedAt        *time.Time `json:"endedAt" form:"endedAt" gorm:"column:ended_at;comment:结束时间"`
	LastMsgID      int64      `json:"lastMsgId" form:"lastMsgId" gorm:"column:last_msg_id;type:bigint;not null;default:0;comment:最后消息ID"`
	LastMsgContent string     `json:"lastMsgContent" form:"lastMsgContent" gorm:"column:last_msg_content;type:varchar(512);not null;default:'';comment:最后消息摘要"`
	LastMsgTime    *time.Time `json:"lastMsgTime" form:"lastMsgTime" gorm:"column:last_msg_time;comment:最后消息时间"`
	UnreadUser     int64      `json:"unreadUser" form:"unreadUser" gorm:"column:unread_user;not null;default:0;comment:用户未读消息数"`
	UnreadClerk    int64      `json:"unreadClerk" form:"unreadClerk" gorm:"column:unread_clerk;not null;default:0;comment:客服未读消息数"`
	Remark         string     `json:"remark" form:"remark" gorm:"column:remark;type:text;comment:备注"`
}

// TableName CsSession 自定义表名 cs_sessions
func (CsSession) TableName() string {
	return "cs_sessions"
}
