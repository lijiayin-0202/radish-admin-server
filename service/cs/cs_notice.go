package cs

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsNoticeService struct{}

// GetCsNoticeList 分页获取站内公告列表，支持标题、关键词、类型与状态筛选
func (s *CsNoticeService) GetCsNoticeList(search csReq.CsNoticeSearch) (list []csModel.CsNotice, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsNotice{})
	if search.Title != "" {
		db = db.Where("title LIKE ?", "%"+search.Title+"%")
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", kw, kw)
	}
	if search.Type != nil {
		db = db.Where("type = ?", *search.Type)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("sort ASC,id DESC").Find(&list).Error
	return
}

// CreateCsNotice 新增站内公告，创建人取当前操作人；类型/状态缺省时取默认值
func (s *CsNoticeService) CreateCsNotice(operatorID uint, req csReq.CsNoticeCreateReq) (entity csModel.CsNotice, err error) {
	if req.Type == 0 {
		req.Type = 1
	}
	if req.Status == 0 {
		req.Status = 1
	}
	entity = csModel.CsNotice{
		Title:     req.Title,
		Content:   req.Content,
		Type:      req.Type,
		Status:    req.Status,
		Sort:      req.Sort,
		CreatedBy: int64(operatorID),
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsNotice 修改站内公告
func (s *CsNoticeService) UpdateCsNotice(operatorID uint, req csReq.CsNoticeUpdateReq) (err error) {
	var entity csModel.CsNotice
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("站内公告不存在或已删除")
	}
	entity.Title = req.Title
	entity.Content = req.Content
	entity.Type = req.Type
	entity.Status = req.Status
	entity.Sort = req.Sort
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsNotice 删除站内公告（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsNoticeService) DeleteCsNotice(id uint) (err error) {
	var entity csModel.CsNotice
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("站内公告不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsNotice 获取站内公告详情
func (s *CsNoticeService) GetCsNotice(id uint) (entity csModel.CsNotice, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
