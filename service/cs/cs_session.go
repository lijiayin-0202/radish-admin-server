package cs

import (
	"errors"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsSessionService struct{}

// 客服会话状态：1进行中，2已结束（取值与 cs_sessions.status 注释一致）
const (
	CsSessionStatusOngoing int64 = 1 // 进行中
	CsSessionStatusEnded   int64 = 2 // 已结束
)

// checkCsSessionStatus 校验会话状态取值
func checkCsSessionStatus(status int64) error {
	if status != CsSessionStatusOngoing && status != CsSessionStatusEnded {
		return errors.New("会话状态取值不合法，仅支持：1进行中，2已结束")
	}
	return nil
}

// GetCsSessionList 分页获取客服会话列表，支持状态、客服、用户与关键词筛选
func (s *CsSessionService) GetCsSessionList(search csReq.CsSessionSearch) (list []csModel.CsSession, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsSession{})
	if search.SessionNo != "" {
		db = db.Where("session_no = ?", search.SessionNo)
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("session_no LIKE ? OR user_name LIKE ? OR clerk_name LIKE ?", kw, kw, kw)
	}
	if search.TicketID != nil {
		db = db.Where("ticket_id = ?", *search.TicketID)
	}
	if search.UserID != nil {
		db = db.Where("user_id = ?", *search.UserID)
	}
	if search.UserPhone != "" {
		db = db.Where("user_phone = ?", search.UserPhone)
	}
	if search.ClerkID != nil {
		db = db.Where("clerk_id = ?", *search.ClerkID)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsSession 新增客服会话
// 会话编号由 global.Snow 雪花ID自动生成，前端不传；新建会话固定为进行中，开始时间为空时取服务器当前时间
// 返回落库后的会话实体，便于前端拿到 sessionId / sessionNo 继续发送会话消息
// TODO(ride-hailing)：会话如需展示订单/车辆信息，调用 ride-hailing 微服务「订单查询接口」获取，
// 本仓库不落地订单/车辆 model
func (s *CsSessionService) CreateCsSession(operatorID uint, req csReq.CsSessionCreateReq) (entity csModel.CsSession, err error) {
	startedAt := req.StartedAt
	if startedAt == nil {
		now := time.Now()
		startedAt = &now
	}
	entity = csModel.CsSession{
		SessionNo: strconv.FormatInt(global.Snow.NextID(), 10), // 会话编号后端自动生成
		UserID:    req.UserID,
		UserName:  req.UserName,
		UserPhone: req.UserPhone,
		ClerkName: req.ClerkName,
		Status:    CsSessionStatusOngoing, // 新建会话固定为进行中
		StartedAt: startedAt,
		Remark:    req.Remark,
	}
	if req.TicketID != nil {
		entity.TicketID = *req.TicketID
	}
	if req.ClerkID != nil {
		entity.ClerkID = *req.ClerkID
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsSession 修改客服会话，会话编号不允许修改
// Status 传 1 进行中 / 2 已结束，nil 表示不变更状态；流转到已结束时结束时间为空则由后端补当前时间
func (s *CsSessionService) UpdateCsSession(operatorID uint, req csReq.CsSessionUpdateReq) (err error) {
	var entity csModel.CsSession
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("会话不存在或已删除")
	}
	if req.Status != nil {
		if err = checkCsSessionStatus(*req.Status); err != nil {
			return
		}
		if *req.Status == CsSessionStatusEnded && req.EndedAt == nil && entity.EndedAt == nil {
			now := time.Now()
			req.EndedAt = &now
		}
		entity.Status = *req.Status
	}

	entity.UserID = req.UserID
	entity.UserName = req.UserName
	entity.UserPhone = req.UserPhone
	entity.ClerkName = req.ClerkName
	entity.Remark = req.Remark
	if req.TicketID != nil {
		entity.TicketID = *req.TicketID
	}
	if req.ClerkID != nil {
		entity.ClerkID = *req.ClerkID
	}
	if req.StartedAt != nil {
		entity.StartedAt = req.StartedAt // 开始时间空值不覆盖已有内容
	}
	if req.EndedAt != nil {
		entity.EndedAt = req.EndedAt // 结束时间空值不覆盖已有内容
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsSession 删除客服会话（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsSessionService) DeleteCsSession(id uint) (err error) {
	var entity csModel.CsSession
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("会话不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsSession 获取客服会话详情
// TODO(ride-hailing)：详情如需展示订单/车辆信息，调用 ride-hailing 微服务「订单查询接口」按 ticket_id 关联工单取 order_id 查询，
// 本服务不落地订单/车辆 model
func (s *CsSessionService) GetCsSession(id uint) (entity csModel.CsSession, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
