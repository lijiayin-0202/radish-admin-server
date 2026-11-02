package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformModel "github.com/flipped-aurora/gin-vue-admin/server/model/platform"
	platformReq "github.com/flipped-aurora/gin-vue-admin/server/model/platform/request"
)

// PricingRuleService 计价规则服务（GVA 主库 pricing_rules）
type PricingRuleService struct{}

func (s *PricingRuleService) CreatePricingRule(a *platformModel.PricingRule) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PricingRuleService) UpdatePricingRule(a *platformModel.PricingRule) error {
	return global.GVA_DB.Model(&platformModel.PricingRule{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PricingRuleService) DeletePricingRule(id uint) error {
	return global.GVA_DB.Delete(&platformModel.PricingRule{}, id).Error
}

func (s *PricingRuleService) DeletePricingRulesByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.PricingRule{}, ids).Error
}

func (s *PricingRuleService) GetPricingRule(id uint) (r platformModel.PricingRule, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&r).Error
	return
}

func (s *PricingRuleService) GetPricingRuleList(info platformReq.PricingRuleSearch) (list []platformModel.PricingRule, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.PricingRule{})
	if info.Keyword != "" {
		db = db.Where("rule_name LIKE ?", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	page, pageSize := info.Page, info.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	err = db.Limit(pageSize).Offset(pageSize * (page - 1)).Order("id desc").Find(&list).Error
	return
}
