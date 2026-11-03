package platform

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformModel "github.com/flipped-aurora/gin-vue-admin/server/model/platform"
	platformReq "github.com/flipped-aurora/gin-vue-admin/server/model/platform/request"
	platformRes "github.com/flipped-aurora/gin-vue-admin/server/model/platform/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"gorm.io/gorm"
)

type PlatformService struct{}

// ============================ 协议 ============================

func (s *PlatformService) GetAgreementList(info platformReq.AgreementSearch) (list []platformModel.Agreement, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.Agreement{})
	if info.Keyword != "" {
		db = db.Where("title LIKE ? OR version LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) GetAgreementByID(id uint) (data platformModel.Agreement, err error) {
	err = global.GVA_DB.First(&data, id).Error
	return
}

func (s *PlatformService) CreateAgreement(a *platformModel.Agreement) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdateAgreement(a *platformModel.Agreement) error {
	return global.GVA_DB.Model(&platformModel.Agreement{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeleteAgreement(id uint) error {
	return global.GVA_DB.Delete(&platformModel.Agreement{}, id).Error
}

func (s *PlatformService) DeleteAgreementsByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.Agreement{}, ids).Error
}

// ============================ 公告 ============================

func (s *PlatformService) GetAnnouncementList(info platformReq.AnnouncementSearch) (list []platformModel.Announcement, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.Announcement{})
	if info.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) GetAnnouncementByID(id uint) (data platformModel.Announcement, err error) {
	err = global.GVA_DB.First(&data, id).Error
	return
}

func (s *PlatformService) CreateAnnouncement(a *platformModel.Announcement) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdateAnnouncement(a *platformModel.Announcement) error {
	return global.GVA_DB.Model(&platformModel.Announcement{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeleteAnnouncement(id uint) error {
	return global.GVA_DB.Delete(&platformModel.Announcement{}, id).Error
}

func (s *PlatformService) DeleteAnnouncementsByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.Announcement{}, ids).Error
}

// ============================ 平台参数 ============================

func (s *PlatformService) GetSettingList(info platformReq.PlatformSettingSearch) (list []platformModel.PlatformSetting, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.PlatformSetting{})
	if info.Keyword != "" {
		db = db.Where("config_key LIKE ? OR description LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id asc").Find(&list).Error
	return
}

func (s *PlatformService) CreateSetting(a *platformModel.PlatformSetting) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdateSetting(a *platformModel.PlatformSetting) error {
	return global.GVA_DB.Model(&platformModel.PlatformSetting{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeleteSetting(id uint) error {
	return global.GVA_DB.Delete(&platformModel.PlatformSetting{}, id).Error
}

func (s *PlatformService) DeleteSettingsByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.PlatformSetting{}, ids).Error
}

// ============================ 角色 ============================

func (s *PlatformService) GetRoleList(info platformReq.RoleSearch) (list []platformModel.Role, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.Role{})
	if info.Keyword != "" {
		db = db.Where("role_name LIKE ? OR role_code LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) GetRoleWithPermissions(id uint) (data platformRes.RoleWithPermissions, err error) {
	var role platformModel.Role
	if err = global.GVA_DB.First(&role, id).Error; err != nil {
		return
	}
	perms, err := s.GetRolePermissions(id)
	if err != nil {
		return
	}
	data = platformRes.RoleWithPermissions{
		ID: role.ID, RoleName: role.RoleName, RoleCode: role.RoleCode,
		Description: role.Description, Status: role.Status, PermissionIDs: perms,
	}
	return
}

func (s *PlatformService) GetRolePermissions(roleID uint) ([]uint, error) {
	var ids []uint
	err := global.GVA_DB.Model(&platformModel.RolePermission{}).Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error
	return ids, err
}

func (s *PlatformService) CreateRole(r *platformModel.Role) error {
	return global.GVA_DB.Create(r).Error
}

func (s *PlatformService) UpdateRole(r *platformModel.Role) error {
	return global.GVA_DB.Model(&platformModel.Role{}).Where("id = ?", r.ID).Omit("created_at").Save(r).Error
}

func (s *PlatformService) DeleteRole(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&platformModel.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&platformModel.Role{}, id).Error
	})
}

func (s *PlatformService) AssignRolePermissions(roleID uint, permIDs []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&platformModel.RolePermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permIDs {
			if err := tx.Create(&platformModel.RolePermission{RoleID: roleID, PermissionID: pid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================ 权限点 ============================

func (s *PlatformService) GetPermissionList(info platformReq.PermissionSearch) (list []platformModel.Permission, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.Permission{})
	if info.Keyword != "" {
		db = db.Where("permission_name LIKE ? OR route_path LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.PermissionType != nil {
		db = db.Where("permission_type = ?", *info.PermissionType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("sort asc, id asc").Find(&list).Error
	return
}

func (s *PlatformService) GetPermissionTree() ([]platformRes.PermissionTree, error) {
	var list []platformModel.Permission
	if err := global.GVA_DB.Order("sort asc, id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return buildPermissionTree(list, 0), nil
}

func buildPermissionTree(list []platformModel.Permission, parentID uint) []platformRes.PermissionTree {
	var nodes []platformRes.PermissionTree
	for _, p := range list {
		if p.ParentID != parentID {
			continue
		}
		node := platformRes.PermissionTree{
			ID: p.ID, ParentID: p.ParentID, PermissionName: p.PermissionName,
			PermissionType: p.PermissionType, RoutePath: p.RoutePath,
			HTTPMethod: p.HTTPMethod, Sort: p.Sort,
		}
		node.Children = buildPermissionTree(list, p.ID)
		nodes = append(nodes, node)
	}
	return nodes
}

func (s *PlatformService) CreatePermission(p *platformModel.Permission) error {
	return global.GVA_DB.Create(p).Error
}

func (s *PlatformService) UpdatePermission(p *platformModel.Permission) error {
	return global.GVA_DB.Model(&platformModel.Permission{}).Where("id = ?", p.ID).Omit("created_at").Save(p).Error
}

func (s *PlatformService) DeletePermission(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("permission_id = ?", id).Delete(&platformModel.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&platformModel.Permission{}, id).Error
	})
}

// ============================ 后台账号 ============================

func (s *PlatformService) GetOperatorList(info platformReq.OperatorUserSearch) (list []platformModel.OperatorUser, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.OperatorUser{})
	if info.Keyword != "" {
		db = db.Where("username LIKE ? OR real_name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) GetOperatorWithRoles(id uint) (data platformRes.OperatorWithRoles, err error) {
	var op platformModel.OperatorUser
	if err = global.GVA_DB.First(&op, id).Error; err != nil {
		return
	}
	var roleIDs []uint
	err = global.GVA_DB.Model(&platformModel.OperatorUserRole{}).Where("operator_user_id = ?", id).Pluck("role_id", &roleIDs).Error
	if err != nil {
		return
	}
	data = platformRes.OperatorWithRoles{
		ID: op.ID, Username: op.Username, RealName: op.RealName, Status: op.Status, RoleIDs: roleIDs,
	}
	return
}

func (s *PlatformService) CreateOperator(u *platformModel.OperatorUser) error {
	u.Password = utils.BcryptHash(u.Password)
	return global.GVA_DB.Create(u).Error
}

func (s *PlatformService) UpdateOperator(u *platformModel.OperatorUser) error {
	if u.Password != "" {
		u.Password = utils.BcryptHash(u.Password)
	}
	return global.GVA_DB.Model(&platformModel.OperatorUser{}).Where("id = ?", u.ID).Omit("created_at").Save(u).Error
}

func (s *PlatformService) DeleteOperator(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("operator_user_id = ?", id).Delete(&platformModel.OperatorUserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&platformModel.OperatorUser{}, id).Error
	})
}

func (s *PlatformService) AssignOperatorRoles(operatorID uint, roleIDs []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("operator_user_id = ?", operatorID).Delete(&platformModel.OperatorUserRole{}).Error; err != nil {
			return err
		}
		for _, rid := range roleIDs {
			if err := tx.Create(&platformModel.OperatorUserRole{OperatorUserID: operatorID, RoleID: rid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ============================ 操作审计 ============================

func (s *PlatformService) GetAuditLogList(info platformReq.AuditLogSearch) (list []platformModel.OperationAuditLog, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.OperationAuditLog{})
	if info.Keyword != "" {
		db = db.Where("content LIKE ? OR operator_name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Module != "" {
		db = db.Where("module = ?", info.Module)
	}
	if info.OperatorName != "" {
		db = db.Where("operator_name LIKE ?", "%"+info.OperatorName+"%")
	}
	if info.StartTime != "" && info.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartTime, info.EndTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) WriteAuditLog(log *platformModel.OperationAuditLog) error {
	return global.GVA_DB.Create(log).Error
}

// ============================ 计价规则 ============================

func (s *PlatformService) GetPricingRuleList(info platformReq.PricingRuleSearch) (list []platformModel.PricingRule, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.PricingRule{})
	if info.Keyword != "" {
		db = db.Where("rule_name LIKE ?", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) CreatePricingRule(a *platformModel.PricingRule) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdatePricingRule(a *platformModel.PricingRule) error {
	return global.GVA_DB.Model(&platformModel.PricingRule{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeletePricingRule(id uint) error {
	return global.GVA_DB.Delete(&platformModel.PricingRule{}, id).Error
}

func (s *PlatformService) DeletePricingRulesByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.PricingRule{}, ids).Error
}

// ============================ 地图版本 ============================

func (s *PlatformService) GetMapVersionList(info platformReq.MapVersionSearch) (list []platformModel.MapVersion, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.MapVersion{})
	if info.Keyword != "" {
		db = db.Where("version LIKE ? OR publisher LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) CreateMapVersion(a *platformModel.MapVersion) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdateMapVersion(a *platformModel.MapVersion) error {
	return global.GVA_DB.Model(&platformModel.MapVersion{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeleteMapVersion(id uint) error {
	return global.GVA_DB.Delete(&platformModel.MapVersion{}, id).Error
}

// ActivateMapVersion 切换生效版本
func (s *PlatformService) ActivateMapVersion(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&platformModel.MapVersion{}).Where("1=1").Update("is_active", false).Error; err != nil {
			return err
		}
		return tx.Model(&platformModel.MapVersion{}).Where("id = ?", id).Update("is_active", true).Error
	})
}

// ============================ 消息推送模板 ============================

func (s *PlatformService) GetPushTemplateList(info platformReq.PushTemplateSearch) (list []platformModel.PushTemplate, total int64, err error) {
	db := global.GVA_DB.Model(&platformModel.PushTemplate{})
	if info.Keyword != "" {
		db = db.Where("template_name LIKE ? OR template_code LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(info.PageSize).Offset(info.PageSize * (info.Page - 1)).Order("id desc").Find(&list).Error
	return
}

func (s *PlatformService) CreatePushTemplate(a *platformModel.PushTemplate) error {
	return global.GVA_DB.Create(a).Error
}

func (s *PlatformService) UpdatePushTemplate(a *platformModel.PushTemplate) error {
	return global.GVA_DB.Model(&platformModel.PushTemplate{}).Where("id = ?", a.ID).Omit("created_at").Save(a).Error
}

func (s *PlatformService) DeletePushTemplate(id uint) error {
	return global.GVA_DB.Delete(&platformModel.PushTemplate{}, id).Error
}

func (s *PlatformService) DeletePushTemplatesByIds(ids []int) error {
	return global.GVA_DB.Delete(&platformModel.PushTemplate{}, ids).Error
}
