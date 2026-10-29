package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsFaq 常见问题表
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
type CsFaq struct {
	global.GVA_MODEL
	Question    string `json:"question" form:"question" gorm:"column:question;type:varchar(512);not null;default:'';comment:问题标题"`
	Answer      string `json:"answer" form:"answer" gorm:"column:answer;type:text;not null;comment:回答内容"`
	Category    string `json:"category" form:"category" gorm:"column:category;type:varchar(128);not null;default:'';comment:问题分类"`
	Tags        string `json:"tags" form:"tags" gorm:"column:tags;type:varchar(256);not null;default:'';comment:标签"`
	Sort        int64  `json:"sort" form:"sort" gorm:"column:sort;not null;default:0;comment:排序号"`
	Status      int64  `json:"status" form:"status" gorm:"column:status;not null;default:1;comment:状态：0禁用，1启用"`
	ViewCount   int64  `json:"viewCount" form:"viewCount" gorm:"column:view_count;type:bigint;not null;default:0;comment:查看次数"`
	UsefulCount int64  `json:"usefulCount" form:"usefulCount" gorm:"column:useful_count;type:bigint;not null;default:0;comment:有用次数"`
	Remark      string `json:"remark" form:"remark" gorm:"column:remark;type:text;comment:备注"`
}

// TableName CsFaq 自定义表名 cs_faqs
func (CsFaq) TableName() string {
	return "cs_faqs"
}
