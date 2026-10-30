package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsFaqCategory 常见问题分类表
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
type CsFaqCategory struct {
	global.GVA_MODEL
	Title  string `json:"title" form:"title" gorm:"column:title;type:varchar(128);not null;default:'';comment:分类名称"`
	Sort   int64  `json:"sort" form:"sort" gorm:"column:sort;not null;default:0;comment:排序号"`
	Status int64  `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:状态：0禁用，1启用"`
	Remark string `json:"remark" form:"remark" gorm:"column:remark;type:text;comment:备注"`
}

// TableName CsFaqCategory 自定义表名 cs_faq_category
func (CsFaqCategory) TableName() string {
	return "cs_faq_category"
}
