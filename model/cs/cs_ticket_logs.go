package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsTicketLog 工单流水(状态变更日志)表
// 记录工单状态流转历史，仅允许新增与查询，不允许修改、删除
// 操作时间即记录创建时间（GVA_MODEL 自带的 CreatedAt），由 DTO 以 createdAt 输出
type CsTicketLog struct {
	global.GVA_MODEL
	TicketID     int64  `json:"ticketId" form:"ticketId" gorm:"column:ticket_id;type:bigint;not null;default:0;comment:关联工单ID"`
	OperatorID   int64  `json:"operatorId" form:"operatorId" gorm:"column:operator_id;type:bigint;not null;default:0;comment:操作人ID"`
	StatusBefore int64  `json:"statusBefore" form:"statusBefore" gorm:"column:status_before;type:bigint;not null;default:0;comment:变更前状态"`
	StatusAfter  int64  `json:"statusAfter" form:"statusAfter" gorm:"column:status_after;type:bigint;not null;default:0;comment:变更后状态"`
	Remark       string `json:"remark" form:"remark" gorm:"column:remark;type:varchar(512);not null;default:'';comment:操作备注"`
}

// TableName CsTicketLog 自定义表名 cs_ticket_logs
func (CsTicketLog) TableName() string {
	return "cs_ticket_logs"
}
