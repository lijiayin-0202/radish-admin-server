package cs

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	"go.uber.org/zap"
)

type CsTicketService struct{}

// 客服工单状态枚举：0待受理，1处理中，2已关闭
// 流转规则（强制校验）：待受理 -> 处理中、待受理 -> 已关闭、处理中 -> 已关闭
// 已关闭为终态，禁止回退为待受理/处理中
const (
	CsTicketStatusPending    int64 = 0 // 待受理
	CsTicketStatusProcessing int64 = 1 // 处理中
	CsTicketStatusClosed     int64 = 2 // 已关闭
)

// csTicketStatusName 工单状态名称，用于错误提示与前端展示对齐
var csTicketStatusName = map[int64]string{
	CsTicketStatusPending:    "待受理",
	CsTicketStatusProcessing: "处理中",
	CsTicketStatusClosed:     "已关闭",
}

// csTicketStatusFlow 允许的工单状态流转，已关闭为终态
var csTicketStatusFlow = map[int64][]int64{
	CsTicketStatusPending:    {CsTicketStatusProcessing, CsTicketStatusClosed},
	CsTicketStatusProcessing: {CsTicketStatusClosed},
}

// checkCsTicketStatus 校验工单状态取值是否合法：0待受理，1处理中，2已关闭
func checkCsTicketStatus(status int64) error {
	if status != CsTicketStatusPending && status != CsTicketStatusProcessing && status != CsTicketStatusClosed {
		return errors.New("工单状态取值不合法，仅支持：0待受理，1处理中，2已关闭")
	}
	return nil
}

// checkCsTicketStatusFlow 校验工单状态流转是否合法：已关闭为终态，禁止回退
func checkCsTicketStatusFlow(oldStatus, newStatus int64) error {
	if oldStatus == newStatus {
		return nil // 状态未变化，不做流转校验
	}
	if oldStatus == CsTicketStatusClosed {
		return errors.New("工单已关闭，处于终态不允许再流转")
	}
	next, ok := csTicketStatusFlow[oldStatus]
	if !ok {
		return errors.New("工单当前状态不合法，无法确定可流转状态")
	}
	for _, s := range next {
		if s == newStatus {
			return nil
		}
	}
	names := make([]string, 0, len(next))
	for _, s := range next {
		names = append(names, csTicketStatusName[s])
	}
	return errors.New("工单状态流转不合法，" + csTicketStatusName[oldStatus] + "只能流转到" + strings.Join(names, "/"))
}

// GetCsTicketList 分页获取客服工单列表，支持工单状态与关键词筛选
func (s *CsTicketService) GetCsTicketList(search csReq.CsTicketSearch) (list []csModel.CsTicket, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsTicket{})
	if search.TicketNo != "" {
		db = db.Where("ticket_no LIKE ?", "%"+search.TicketNo+"%")
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("ticket_no LIKE ? OR title LIKE ? OR description LIKE ?", kw, kw, kw)
	}
	if search.UserID != nil {
		db = db.Where("user_id = ?", *search.UserID)
	}
	if search.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+search.UserName+"%")
	}
	if search.UserPhone != "" {
		db = db.Where("user_phone = ?", search.UserPhone)
	}
	if search.OrderID != nil {
		db = db.Where("order_id = ?", *search.OrderID)
	}
	if search.OrderNo != "" {
		db = db.Where("order_no = ?", search.OrderNo)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	if search.Priority != nil {
		db = db.Where("priority = ?", *search.Priority)
	}
	if search.Category != "" {
		db = db.Where("category LIKE ?", "%"+search.Category+"%")
	}
	if search.AssignedClerkID != nil {
		db = db.Where("assigned_clerk_id = ?", *search.AssignedClerkID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsTicket 新增客服工单
// 工单号由 global.Snow 雪花ID自动生成，前端不传；新建工单固定为待受理
func (s *CsTicketService) CreateCsTicket(operatorID uint, req csReq.CsTicketCreateReq) (entity csModel.CsTicket, err error) {
	// 优先级取值：1低，2中，3高（此处不能用状态常量，状态枚举已改为 0/1/2）
	if req.Priority < 1 || req.Priority > 3 {
		return entity, errors.New("工单优先级取值不合法，仅支持：1低，2中，3高")
	}
	entity = csModel.CsTicket{
		TicketNo:          strconv.FormatInt(global.Snow.NextID(), 10), // 工单号后端自动生成
		UserID:            req.UserID,
		UserName:          req.UserName,
		UserPhone:         req.UserPhone,
		OrderNo:           req.OrderNo,
		Priority:          req.Priority,
		Category:          req.Category,
		Title:             req.Title,
		Description:       req.Description,
		AttachUrls:        req.AttachUrls,
		AssignedClerkName: req.AssignedClerkName,
		Status:            CsTicketStatusPending, // 新建工单固定为待受理
		CreatedBy:         int64(operatorID),     // 创建者：取当前操作人ID
	}
	if req.OrderID != nil {
		entity.OrderID = *req.OrderID
	}
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}

	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsTicket 修改客服工单，工单号不允许修改
// 状态变更按流转规则校验：待受理 -> 处理中 / 待受理 -> 已关闭 / 处理中 -> 已关闭；req.Status 为 nil 表示不变更状态
func (s *CsTicketService) UpdateCsTicket(operatorID uint, req csReq.CsTicketUpdateReq) (err error) {
	var entity csModel.CsTicket
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("工单不存在或已删除")
	}
	if req.Priority != 0 && (req.Priority < 1 || req.Priority > 3) {
		return errors.New("工单优先级取值不合法，仅支持：1低，2中，3高")
	}

	// 状态流转校验
	if req.Status != nil {
		if err = checkCsTicketStatus(*req.Status); err != nil {
			return
		}
		if err = checkCsTicketStatusFlow(entity.Status, *req.Status); err != nil {
			return
		}
		// 流转到处理中，必须指定处理客服
		if *req.Status == CsTicketStatusProcessing && req.AssignedClerkID == nil && entity.AssignedClerkID == 0 {
			return errors.New("工单流转到处理中，必须指定处理客服")
		}
		// 流转到已关闭，自动补解决时间与关闭时间
		if *req.Status == CsTicketStatusClosed {
			now := time.Now()
			if req.ResolvedAt == nil && entity.ResolvedAt == nil {
				req.ResolvedAt = &now
			}
			if req.ClosedAt == nil && entity.ClosedAt == nil {
				req.ClosedAt = &now
			}
		}
		entity.Status = *req.Status
	}

	entity.UserID = req.UserID
	entity.UserName = req.UserName
	entity.UserPhone = req.UserPhone
	entity.OrderNo = req.OrderNo
	entity.Category = req.Category
	entity.Title = req.Title
	entity.Description = req.Description
	entity.AttachUrls = req.AttachUrls
	entity.AssignedClerkName = req.AssignedClerkName
	entity.Remark = req.Remark
	if req.OrderID != nil {
		entity.OrderID = *req.OrderID
	}
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}
	if req.Priority != 0 {
		entity.Priority = req.Priority // 优先级为0表示不变更
	}
	if req.ResolvedAt != nil {
		entity.ResolvedAt = req.ResolvedAt // 解决时间空值不覆盖已有内容
	}
	if req.ClosedAt != nil {
		entity.ClosedAt = req.ClosedAt // 关闭时间空值不覆盖已有内容
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsTicket 删除客服工单（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsTicketService) DeleteCsTicket(operatorID uint, id uint) (err error) {
	var entity csModel.CsTicket
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("工单不存在或已删除")
	}
	global.GVA_LOG.Info("删除客服工单", zap.Uint("operatorId", operatorID), zap.Uint("ticketId", id))
	return global.LUOBO_DB.Delete(&entity).Error
}

// UpdateCsTicketStatus 修改客服工单状态
// 必填入参只有工单ID与状态，处理客服、解决/关闭时间、备注均为可选补充字段
// 状态流转强制校验：待受理(0) -> 处理中(1)、待受理(0) -> 已关闭(2)、处理中(1) -> 已关闭(2)；
// 已关闭(2)为终态，禁止回退到待受理/处理中；状态取值或流转非法时返回业务错误，不 panic
func (s *CsTicketService) UpdateCsTicketStatus(operatorID uint, req csReq.CsTicketStatusReq) (entity csModel.CsTicket, err error) {
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return entity, errors.New("工单不存在或已删除")
	}
	// 目标状态取值校验
	if err = checkCsTicketStatus(req.Status); err != nil {
		return
	}
	// 状态流转校验：禁止乱跳、禁止已关闭回退
	oldStatus := entity.Status
	if err = checkCsTicketStatusFlow(oldStatus, req.Status); err != nil {
		return
	}

	now := time.Now()
	// 流转到已关闭，自动补关闭时间；未填解决时间时按当前时间一并补齐
	if req.Status == CsTicketStatusClosed {
		if req.ResolvedAt == nil && entity.ResolvedAt == nil {
			req.ResolvedAt = &now
		}
		if req.ClosedAt == nil && entity.ClosedAt == nil {
			req.ClosedAt = &now
		}
	}

	entity.Status = req.Status
	if req.AssignedClerkID != nil {
		entity.AssignedClerkID = *req.AssignedClerkID
	}
	if req.AssignedClerkName != "" {
		entity.AssignedClerkName = req.AssignedClerkName
	}
	if req.ResolvedAt != nil {
		entity.ResolvedAt = req.ResolvedAt // 解决时间空值不覆盖已有内容
	}
	if req.ClosedAt != nil {
		entity.ClosedAt = req.ClosedAt // 关闭时间空值不覆盖已有内容
	}
	if req.Remark != "" {
		entity.Remark = req.Remark
	}
	entity.UpdatedAt = now // 状态变更刷新更新时间

	// 状态变更记录操作人（operatorID 来自 api 层 utils.GetUserID(c)，无任何兜底）
	global.GVA_LOG.Info("修改客服工单状态",
		zap.Uint("operatorId", operatorID),
		zap.Uint("ticketId", req.ID),
		zap.Int64("oldStatus", oldStatus),
		zap.Int64("newStatus", req.Status))

	err = global.LUOBO_DB.Save(&entity).Error
	return
}

// GetCsTicket 获取客服工单详情
func (s *CsTicketService) GetCsTicket(operatorID uint, id uint) (entity csModel.CsTicket, err error) {
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, errors.New("工单不存在或已删除")
	}
	global.GVA_LOG.Info("查询客服工单详情", zap.Uint("operatorId", operatorID), zap.Uint("ticketId", id))
	return
}
