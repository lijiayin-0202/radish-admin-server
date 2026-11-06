package cs

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsTicketLogService struct{}

// GetCsTicketLogList 分页获取工单流水列表
// 支持按工单ID、操作人、变更后状态与关键词（操作备注）筛选，按主键倒序（最新操作在前）
func (s *CsTicketLogService) GetCsTicketLogList(search csReq.CsTicketLogSearch) (list []csModel.CsTicketLog, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsTicketLog{})
	if search.TicketID != nil {
		db = db.Where("ticket_id = ?", *search.TicketID)
	}
	if search.OperatorID != nil {
		db = db.Where("operator_id = ?", *search.OperatorID)
	}
	if search.StatusAfter != nil {
		db = db.Where("status_after = ?", *search.StatusAfter)
	}
	if search.Keyword != "" {
		db = db.Where("remark LIKE ?", "%"+search.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsTicketLog 新增工单流水
// 操作人ID由调用方（api 层当前登录用户）传入；操作时间取记录创建时间
func (s *CsTicketLogService) CreateCsTicketLog(operatorID uint, req csReq.CsTicketLogCreateReq) (entity csModel.CsTicketLog, err error) {
	if req.TicketID <= 0 {
		return entity, errors.New("关联工单ID不合法")
	}
	entity = csModel.CsTicketLog{
		TicketID:     req.TicketID,
		OperatorID:   int64(operatorID),
		StatusBefore: req.StatusBefore,
		StatusAfter:  req.StatusAfter,
		Remark:       req.Remark,
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// GetCsTicketLog 获取工单流水详情
func (s *CsTicketLogService) GetCsTicketLog(id uint) (entity csModel.CsTicketLog, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
