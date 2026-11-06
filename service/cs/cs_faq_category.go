package cs

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsFaqCategoryService struct{}

// GetCsFaqCategoryList 分页获取常见问题分类列表，支持分类名称与启用状态筛选
func (s *CsFaqCategoryService) GetCsFaqCategoryList(search csReq.CsFaqCategorySearch) (list []csModel.CsFaqCategory, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsFaqCategory{})
	if search.Title != "" {
		db = db.Where("title LIKE ?", "%"+search.Title+"%")
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("sort DESC").Order("id DESC").Find(&list).Error
	return
}

// CreateCsFaqCategory 新增常见问题分类
// Status 传 0 禁用 / 1 启用，nil 表示按库默认1启用
func (s *CsFaqCategoryService) CreateCsFaqCategory(operatorID uint, req csReq.CsFaqCategoryCreateReq) (entity csModel.CsFaqCategory, err error) {
	status := req.Status
	if status == nil {
		on := CsFaqStatusOn
		status = &on
	}
	if err = checkCsFaqStatus(*status); err != nil {
		return entity, err
	}
	entity = csModel.CsFaqCategory{
		Title:  req.Title,
		Sort:   req.Sort,
		Status: *status,
		Remark: req.Remark,
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsFaqCategory 修改常见问题分类，含启停：Status 传 0 禁用 / 1 启用，nil 表示不变更状态
func (s *CsFaqCategoryService) UpdateCsFaqCategory(operatorID uint, req csReq.CsFaqCategoryUpdateReq) (err error) {
	var entity csModel.CsFaqCategory
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("FAQ分类不存在或已删除")
	}
	entity.Title = req.Title
	entity.Sort = req.Sort
	entity.Remark = req.Remark
	if req.Status != nil {
		if err = checkCsFaqStatus(*req.Status); err != nil {
			return
		}
		entity.Status = *req.Status
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsFaqCategory 删除常见问题分类（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsFaqCategoryService) DeleteCsFaqCategory(id uint) (err error) {
	var entity csModel.CsFaqCategory
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("FAQ分类不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsFaqCategory 获取常见问题分类详情
func (s *CsFaqCategoryService) GetCsFaqCategory(id uint) (entity csModel.CsFaqCategory, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
