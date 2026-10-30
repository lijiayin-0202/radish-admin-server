package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsNotice 消息通知表
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
type CsNotice struct {
	global.GVA_MODEL
	UserID  int64  `json:"userId" form:"userId" gorm:"column:user_id;type:bigint;not null;default:0;comment:接收人用户ID"`
	Title   string `json:"title" form:"title" gorm:"column:title;type:varchar(256);not null;default:'';comment:消息标题"`
	Content string `json:"content" form:"content" gorm:"column:content;type:text;not null;comment:消息内容"`
	IsRead  int64  `json:"isRead" form:"isRead" gorm:"column:is_read;not null;default:0;comment:是否已读：0未读，1已读"`
	Type    int64  `json:"type" form:"type" gorm:"column:type;not null;default:1;comment:消息类型：1系统通知，2工单通知，3会话通知"`
}

// TableName CsNotice 自定义表名 cs_notice
func (CsNotice) TableName() string {
	return "cs_notice"
}
