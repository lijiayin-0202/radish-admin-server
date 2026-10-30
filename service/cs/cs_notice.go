package cs

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsNoticeService struct{}

// 消息通知已读状态：0未读，1已读（取值与 cs_notice.is_read 注释一致）
const (
	CsNoticeUnread int64 = 0 // 未读
	CsNoticeRead   int64 = 1 // 已读
)

// 消息通知类型：1系统通知，2工单通知，3会话通知（取值与 cs_notice.type 注释一致）
const (
	CsNoticeTypeSystem  int64 = 1 // 系统通知
	CsNoticeTypeTicket  int64 = 2 // 工单通知
	CsNoticeTypeSession int64 = 3 // 会话通知
)

// checkCsNoticeType 校验消息通知类型取值
func checkCsNoticeType(noticeType int64) error {
	if noticeType != CsNoticeTypeSystem && noticeType != CsNoticeTypeTicket && noticeType != CsNoticeTypeSession {
		return errors.New("消息类型取值不合法，仅支持：1系统通知，2工单通知，3会话通知")
	}
	return nil
}

// GetCsNoticeList 分页获取消息通知列表，支持接收人、已读状态、消息类型与关键词筛选
func (s *CsNoticeService) GetCsNoticeList(search csReq.CsNoticeSearch) (list []csModel.CsNotice, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsNotice{})
	if search.UserID != nil {
		db = db.Where("user_id = ?", *search.UserID)
	}
	if search.Title != "" {
		db = db.Where("title LIKE ?", "%"+search.Title+"%")
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", kw, kw)
	}
	if search.IsRead != nil {
		db = db.Where("is_read = ?", *search.IsRead)
	}
	if search.Type != nil {
		db = db.Where("type = ?", *search.Type)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsNotice 新增消息通知，新增固定为未读，已读状态由 ReadCsNotice 标记
func (s *CsNoticeService) CreateCsNotice(operatorID uint, req csReq.CsNoticeCreateReq) (entity csModel.CsNotice, err error) {
	if err = checkCsNoticeType(req.Type); err != nil {
		return entity, err
	}
	entity = csModel.CsNotice{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		IsRead:  CsNoticeUnread, // 新增消息固定未读
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsNotice 修改消息通知，已读状态建议走 /read 接口标记
// UserID / Type / IsRead 使用指针：传值表示变更，nil 表示不变更
func (s *CsNoticeService) UpdateCsNotice(operatorID uint, req csReq.CsNoticeUpdateReq) (err error) {
	var entity csModel.CsNotice
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("消息通知不存在或已删除")
	}
	entity.Title = req.Title
	entity.Content = req.Content
	if req.UserID != nil {
		entity.UserID = *req.UserID
	}
	if req.Type != nil {
		if err = checkCsNoticeType(*req.Type); err != nil {
			return
		}
		entity.Type = *req.Type
	}
	if req.IsRead != nil {
		if *req.IsRead != CsNoticeUnread && *req.IsRead != CsNoticeRead {
			return errors.New("已读状态取值不合法，仅支持：0未读，1已读")
		}
		entity.IsRead = *req.IsRead
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// ReadCsNotice 标记消息通知已读，重复标记幂等返回，不做报错
func (s *CsNoticeService) ReadCsNotice(id uint) (err error) {
	var entity csModel.CsNotice
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("消息通知不存在或已删除")
	}
	if entity.IsRead == CsNoticeRead {
		return nil // 已读则直接返回，避免重复写库
	}
	entity.IsRead = CsNoticeRead
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsNotice 删除消息通知（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsNoticeService) DeleteCsNotice(id uint) (err error) {
	var entity csModel.CsNotice
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("消息通知不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsNotice 获取消息通知详情
func (s *CsNoticeService) GetCsNotice(id uint) (entity csModel.CsNotice, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
