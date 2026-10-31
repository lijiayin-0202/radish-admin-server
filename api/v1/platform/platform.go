package platform

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	platformModel "github.com/flipped-aurora/gin-vue-admin/server/model/platform"
	platformReq "github.com/flipped-aurora/gin-vue-admin/server/model/platform/request"
	platformService "github.com/flipped-aurora/gin-vue-admin/server/service/platform"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type PlatformApi struct{}

func strToUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

// writeAudit 写入操作审计日志（失败不阻断主流程）
func writeAudit(c *gin.Context, module, action, content, result string) {
	_ = platformService.PlatformServiceApp.WriteAuditLog(&platformModel.OperationAuditLog{
		OperatorID:   utils.GetUserID(c),
		OperatorName: utils.GetUserName(c),
		Module:       module,
		ActionType:   action,
		Content:      content,
		Result:       result,
		SourceIP:     c.ClientIP(),
	})
}

// ============================ 协议 ============================

func (a *PlatformApi) GetAgreementList(c *gin.Context) {
	var info platformReq.AgreementSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetAgreementList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) GetAgreementById(c *gin.Context) {
	id := c.Param("id")
	data, err := platformService.PlatformServiceApp.GetAgreementByID(strToUint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *PlatformApi) CreateAgreement(c *gin.Context) {
	var req platformModel.Agreement
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreateAgreement(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "协议管理", "新增", "新增协议: "+req.Title, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateAgreement(c *gin.Context) {
	var req platformModel.Agreement
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdateAgreement(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "协议管理", "编辑", "编辑协议: "+req.Title, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteAgreement(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteAgreement(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "协议管理", "删除", "删除协议ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) DeleteAgreementsByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idArr := make([]int, 0, len(ids))
	for _, v := range ids {
		idArr = append(idArr, int(strToUint(v)))
	}
	if err := platformService.PlatformServiceApp.DeleteAgreementsByIds(idArr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "协议管理", "批量删除", "批量删除协议", "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 公告 ============================

func (a *PlatformApi) GetAnnouncementList(c *gin.Context) {
	var info platformReq.AnnouncementSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetAnnouncementList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) CreateAnnouncement(c *gin.Context) {
	var req platformModel.Announcement
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreateAnnouncement(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "公告管理", "新增", "新增公告: "+req.Title, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateAnnouncement(c *gin.Context) {
	var req platformModel.Announcement
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdateAnnouncement(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "公告管理", "编辑", "编辑公告: "+req.Title, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteAnnouncement(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteAnnouncement(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "公告管理", "删除", "删除公告ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) DeleteAnnouncementsByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idArr := make([]int, 0, len(ids))
	for _, v := range ids {
		idArr = append(idArr, int(strToUint(v)))
	}
	if err := platformService.PlatformServiceApp.DeleteAnnouncementsByIds(idArr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "公告管理", "批量删除", "批量删除公告", "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 平台参数 ============================

func (a *PlatformApi) GetSettingList(c *gin.Context) {
	var info platformReq.PlatformSettingSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetSettingList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) CreateSetting(c *gin.Context) {
	var req platformModel.PlatformSetting
	_ = c.ShouldBindJSON(&req)
	req.UpdatedBy = utils.GetUserName(c)
	if err := platformService.PlatformServiceApp.CreateSetting(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "平台参数", "新增", "新增配置: "+req.ConfigKey, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateSetting(c *gin.Context) {
	var req platformModel.PlatformSetting
	_ = c.ShouldBindJSON(&req)
	req.UpdatedBy = utils.GetUserName(c)
	if err := platformService.PlatformServiceApp.UpdateSetting(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "平台参数", "编辑", "修改配置: "+req.ConfigKey, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteSetting(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteSetting(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "平台参数", "删除", "删除配置ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) DeleteSettingsByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idArr := make([]int, 0, len(ids))
	for _, v := range ids {
		idArr = append(idArr, int(strToUint(v)))
	}
	if err := platformService.PlatformServiceApp.DeleteSettingsByIds(idArr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "平台参数", "批量删除", "批量删除配置", "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 角色 ============================

func (a *PlatformApi) GetRoleList(c *gin.Context) {
	var info platformReq.RoleSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetRoleList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) GetRolePermissions(c *gin.Context) {
	id := strToUint(c.Param("id"))
	data, err := platformService.PlatformServiceApp.GetRoleWithPermissions(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *PlatformApi) CreateRole(c *gin.Context) {
	var req platformModel.Role
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreateRole(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "角色管理", "新增", "新增角色: "+req.RoleName, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateRole(c *gin.Context) {
	var req platformModel.Role
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdateRole(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "角色管理", "编辑", "编辑角色: "+req.RoleName, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteRole(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteRole(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "角色管理", "删除", "删除角色ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) AssignRolePermissions(c *gin.Context) {
	var req platformReq.RolePermissionReq
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.AssignRolePermissions(req.RoleID, req.PermissionIDs); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "角色管理", "分配权限", "角色ID: "+strconv.FormatUint(uint64(req.RoleID), 10)+" 分配权限", "成功")
	response.OkWithMessage("分配成功", c)
}

// ============================ 权限点 ============================

func (a *PlatformApi) GetPermissionList(c *gin.Context) {
	var info platformReq.PermissionSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetPermissionList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) GetPermissionTree(c *gin.Context) {
	tree, err := platformService.PlatformServiceApp.GetPermissionTree()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(tree, "获取成功", c)
}

func (a *PlatformApi) CreatePermission(c *gin.Context) {
	var req platformModel.Permission
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreatePermission(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "权限管理", "新增", "新增权限: "+req.PermissionName, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdatePermission(c *gin.Context) {
	var req platformModel.Permission
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdatePermission(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "权限管理", "编辑", "编辑权限: "+req.PermissionName, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeletePermission(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeletePermission(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "权限管理", "删除", "删除权限ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 后台账号 ============================

func (a *PlatformApi) GetOperatorList(c *gin.Context) {
	var info platformReq.OperatorUserSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetOperatorList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) GetOperatorRoles(c *gin.Context) {
	id := strToUint(c.Param("id"))
	data, err := platformService.PlatformServiceApp.GetOperatorWithRoles(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *PlatformApi) CreateOperator(c *gin.Context) {
	var req platformModel.OperatorUser
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreateOperator(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "账号管理", "新增", "新增账号: "+req.Username, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateOperator(c *gin.Context) {
	var req platformModel.OperatorUser
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdateOperator(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "账号管理", "编辑", "编辑账号: "+req.Username, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteOperator(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteOperator(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "账号管理", "删除", "删除账号ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) AssignOperatorRoles(c *gin.Context) {
	var req platformReq.OperatorUserRoleReq
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.AssignOperatorRoles(req.OperatorUserID, req.RoleIDs); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "账号管理", "分配角色", "账号ID: "+strconv.FormatUint(uint64(req.OperatorUserID), 10)+" 分配角色", "成功")
	response.OkWithMessage("分配成功", c)
}

// ============================ 审计日志 ============================

func (a *PlatformApi) GetAuditLogList(c *gin.Context) {
	var info platformReq.AuditLogSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetAuditLogList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// ============================ 计价规则 ============================

func (a *PlatformApi) GetPricingRuleList(c *gin.Context) {
	var info platformReq.PricingRuleSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetPricingRuleList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) CreatePricingRule(c *gin.Context) {
	var req platformModel.PricingRule
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreatePricingRule(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "计价规则", "新增", "新增计价规则: "+req.RuleName, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdatePricingRule(c *gin.Context) {
	var req platformModel.PricingRule
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdatePricingRule(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "计价规则", "编辑", "编辑计价规则: "+req.RuleName, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeletePricingRule(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeletePricingRule(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "计价规则", "删除", "删除计价规则ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) DeletePricingRulesByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idArr := make([]int, 0, len(ids))
	for _, v := range ids {
		idArr = append(idArr, int(strToUint(v)))
	}
	if err := platformService.PlatformServiceApp.DeletePricingRulesByIds(idArr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "计价规则", "批量删除", "批量删除计价规则", "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 地图版本 ============================

func (a *PlatformApi) GetMapVersionList(c *gin.Context) {
	var info platformReq.MapVersionSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetMapVersionList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) CreateMapVersion(c *gin.Context) {
	var req platformModel.MapVersion
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreateMapVersion(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "地图版本", "新增", "新增地图版本: "+req.Version, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdateMapVersion(c *gin.Context) {
	var req platformModel.MapVersion
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdateMapVersion(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "地图版本", "编辑", "编辑地图版本: "+req.Version, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeleteMapVersion(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeleteMapVersion(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "地图版本", "删除", "删除地图版本ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) ActivateMapVersion(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.ActivateMapVersion(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "地图版本", "切换生效", "切换生效地图版本ID: "+c.Param("id"), "成功")
	response.OkWithMessage("切换成功", c)
}

// ============================ 消息推送模板 ============================

func (a *PlatformApi) GetPushTemplateList(c *gin.Context) {
	var info platformReq.PushTemplateSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := platformService.PlatformServiceApp.GetPushTemplateList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *PlatformApi) CreatePushTemplate(c *gin.Context) {
	var req platformModel.PushTemplate
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.CreatePushTemplate(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "消息模板", "新增", "新增消息模板: "+req.TemplateName, "成功")
	response.OkWithMessage("创建成功", c)
}

func (a *PlatformApi) UpdatePushTemplate(c *gin.Context) {
	var req platformModel.PushTemplate
	_ = c.ShouldBindJSON(&req)
	if err := platformService.PlatformServiceApp.UpdatePushTemplate(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "消息模板", "编辑", "编辑消息模板: "+req.TemplateName, "成功")
	response.OkWithMessage("更新成功", c)
}

func (a *PlatformApi) DeletePushTemplate(c *gin.Context) {
	id := strToUint(c.Param("id"))
	if err := platformService.PlatformServiceApp.DeletePushTemplate(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "消息模板", "删除", "删除消息模板ID: "+c.Param("id"), "成功")
	response.OkWithMessage("删除成功", c)
}

func (a *PlatformApi) DeletePushTemplatesByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idArr := make([]int, 0, len(ids))
	for _, v := range ids {
		idArr = append(idArr, int(strToUint(v)))
	}
	if err := platformService.PlatformServiceApp.DeletePushTemplatesByIds(idArr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	writeAudit(c, "消息模板", "批量删除", "批量删除消息模板", "成功")
	response.OkWithMessage("删除成功", c)
}

// ============================ 统计大屏 ============================

func (a *PlatformApi) GetDashboardOverview(c *gin.Context) {
	var q platformReq.DashboardQuery
	_ = c.ShouldBindQuery(&q)
	data, err := platformService.PlatformServiceApp.GetDashboardOverview(q)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *PlatformApi) GetVehicleHeatmap(c *gin.Context) {
	data, err := platformService.PlatformServiceApp.GetVehicleHeatmap()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *PlatformApi) GetStationFlow(c *gin.Context) {
	var q platformReq.DashboardQuery
	_ = c.ShouldBindQuery(&q)
	data, err := platformService.PlatformServiceApp.GetStationFlow(q)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetDashboardCities 大屏城市下拉选项
func (a *PlatformApi) GetDashboardCities(c *gin.Context) {
	data, err := platformService.PlatformServiceApp.GetCities()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetDashboardAreas 大屏区域下拉选项（按城市联动）
func (a *PlatformApi) GetDashboardAreas(c *gin.Context) {
	data, err := platformService.PlatformServiceApp.GetAreas(c.Query("city"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}
