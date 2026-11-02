package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsNotice 站内公告表（后台管理增删改查，与 cs_notifications 用户通知推送相互独立）
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
type CsNotice struct {
	global.GVA_MODEL
	Title     string `json:"title" form:"title" gorm:"column:title;type:varchar(256);not null;default:'';comment:公告标题"`
	Content   string `json:"content" form:"content" gorm:"column:content;type:text;not null;comment:公告内容"`
	Type      int64  `json:"type" form:"type" gorm:"column:type;type:int;not null;default:1;comment:公告类型：1系统公告，2活动公告，3维护公告"`
	Status    int64  `json:"status" form:"status" gorm:"column:status;type:int;not null;default:1;comment:状态：0草稿，1已发布"`
	Sort      int64  `json:"sort" form:"sort" gorm:"column:sort;type:int;not null;default:0;comment:排序号"`
	CreatedBy int64  `json:"createdBy" form:"createdBy" gorm:"column:created_by;type:bigint;not null;default:0;comment:创建者（操作人ID）"`
}

// TableName CsNotice 自定义表名 cs_notice
func (CsNotice) TableName() string {
	return "cs_notice"
}
