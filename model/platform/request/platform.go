package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// PageInfo 通用分页 + 关键词
// AgreementSearch 协议查询
type AgreementSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// AnnouncementSearch 公告查询
type AnnouncementSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// PlatformSettingSearch 平台参数查询
type PlatformSettingSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

// RoleSearch 角色查询
type RoleSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// PermissionSearch 权限点查询
type PermissionSearch struct {
	request.PageInfo
	Keyword        string `json:"keyword" form:"keyword"`
	PermissionType *int   `json:"permissionType" form:"permissionType"`
}

// OperatorUserSearch 后台账号查询
type OperatorUserSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

// AuditLogSearch 审计日志查询
type AuditLogSearch struct {
	request.PageInfo
	Keyword     string `json:"keyword" form:"keyword"`
	Module      string `json:"module" form:"module"`
	OperatorName string `json:"operatorName" form:"operatorName"`
	StartTime   string `json:"startTime" form:"startTime"`
	EndTime     string `json:"endTime" form:"endTime"`
}

// PricingRuleSearch 计价规则查询
type PricingRuleSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
	Status  *int   `json:"status" form:"status"`
}

// MapVersionSearch 地图版本查询
type MapVersionSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

// PushTemplateSearch 消息模板查询
type PushTemplateSearch struct {
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

// RolePermissionReq 角色分配权限
type RolePermissionReq struct {
	RoleID        uint   `json:"roleId"`
	PermissionIDs []uint `json:"permissionIds"`
}

// OperatorUserRoleReq 账号分配角色
type OperatorUserRoleReq struct {
	OperatorUserID uint   `json:"operatorUserId"`
	RoleIDs        []uint `json:"roleIds"`
}

// DashboardQuery 大屏查询（支持时间 + 城市/区域筛选）
type DashboardQuery struct {
	StartTime string `json:"startTime" form:"startTime"` // 开始时间 YYYY-MM-DD HH:MM:SS
	EndTime   string `json:"endTime" form:"endTime"`     // 结束时间
	City      string `json:"city" form:"city"`           // 城市编码（空=全部城市）
	Area      uint   `json:"area" form:"area"`           // 运营区域ID（空=该城市全部区域）
}
