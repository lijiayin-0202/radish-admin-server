package platform

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ==================== 协议 ====================

// Agreement 平台协议（版本化管理）
type Agreement struct {
	global.GVA_MODEL
	Version     string     `json:"version" gorm:"type:varchar(50);comment:协议版本号"`
	Title       string     `json:"title" gorm:"type:varchar(200);comment:协议标题"`
	Content     string     `json:"content" gorm:"type:text;comment:协议内容"`
	Status      int        `json:"status" gorm:"type:int;default:0;comment:状态 0草稿 1生效 2失效"`
	EffectiveAt *time.Time `json:"effectiveAt" gorm:"comment:生效时间"`
}

func (Agreement) TableName() string { return "agreements" }

// ==================== 公告 ====================

// Announcement 运营公告
type Announcement struct {
	global.GVA_MODEL
	Title     string     `json:"title" gorm:"type:varchar(200);comment:公告标题"`
	Content   string     `json:"content" gorm:"type:text;comment:公告内容"`
	Region    string     `json:"region" gorm:"type:varchar(200);comment:区域范围"`
	Status    int        `json:"status" gorm:"type:int;default:0;comment:状态 0草稿 1已发布 2已下线"`
	PublishAt *time.Time `json:"publishAt" gorm:"comment:发布时间"`
	Publisher string     `json:"publisher" gorm:"type:varchar(100);comment:发布人"`
}

func (Announcement) TableName() string { return "announcements" }

// ==================== 平台参数 ====================

// PlatformSetting 平台参数（键值配置）
type PlatformSetting struct {
	global.GVA_MODEL
	ConfigKey   string `json:"configKey" gorm:"type:varchar(100);uniqueIndex;comment:配置键"`
	ConfigValue string `json:"configValue" gorm:"type:text;comment:配置值"`
	Description string `json:"description" gorm:"type:varchar(255);comment:配置说明"`
	UpdatedBy   string `json:"updatedBy" gorm:"type:varchar(100);comment:最后更新人"`
}

func (PlatformSetting) TableName() string { return "platform_settings" }

// ==================== RBAC ====================

// Role 后台角色
type Role struct {
	global.GVA_MODEL
	RoleName    string `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色名称"`
	RoleCode    string `json:"roleCode" gorm:"type:varchar(100);comment:角色编码"`
	Description string `json:"description" gorm:"type:varchar(255);comment:角色描述"`
	Status      int    `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
}

func (Role) TableName() string { return "roles" }

// Permission 权限点（菜单/按钮/接口）
type Permission struct {
	global.GVA_MODEL
	ParentID       uint   `json:"parentId" gorm:"comment:父级权限ID 0为顶层"`
	PermissionName string `json:"permissionName" gorm:"type:varchar(100);comment:权限名称"`
	PermissionType int    `json:"permissionType" gorm:"type:int;comment:类型 1菜单 2按钮 3接口"`
	RoutePath      string `json:"routePath" gorm:"type:varchar(200);comment:前端路由或接口路径"`
	HTTPMethod     string `json:"httpMethod" gorm:"type:varchar(10);comment:接口HTTP方法"`
	Sort           int    `json:"sort" gorm:"type:int;default:0;comment:排序"`
}

func (Permission) TableName() string { return "permissions" }

// RolePermission 角色-权限关联
type RolePermission struct {
	global.GVA_MODEL
	RoleID       uint `json:"roleId" gorm:"comment:角色ID"`
	PermissionID uint `json:"permissionId" gorm:"comment:权限ID"`
}

func (RolePermission) TableName() string { return "role_permissions" }

// OperatorUser 后台账号
type OperatorUser struct {
	global.GVA_MODEL
	Username string `json:"username" gorm:"type:varchar(100);uniqueIndex;comment:登录账号"`
	Password string `json:"-" gorm:"type:varchar(200);comment:密码哈希"`
	RealName string `json:"realName" gorm:"type:varchar(100);comment:真实姓名"`
	Status   int    `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
}

func (OperatorUser) TableName() string { return "operator_users" }

// OperatorUserRole 账号-角色关联
type OperatorUserRole struct {
	global.GVA_MODEL
	OperatorUserID uint `json:"operatorUserId" gorm:"comment:后台账号ID"`
	RoleID         uint `json:"roleId" gorm:"comment:角色ID"`
}

func (OperatorUserRole) TableName() string { return "operator_user_roles" }

// ==================== 操作审计 ====================

// OperationAuditLog 操作审计日志（只读）
type OperationAuditLog struct {
	global.GVA_MODEL
	OperatorID    uint   `json:"operatorId" gorm:"comment:操作人ID"`
	OperatorName  string `json:"operatorName" gorm:"type:varchar(100);comment:操作人"`
	Module        string `json:"module" gorm:"type:varchar(100);comment:操作模块"`
	ActionType    string `json:"actionType" gorm:"type:varchar(50);comment:操作类型"`
	Content       string `json:"content" gorm:"type:varchar(255);comment:操作内容"`
	RequestParams string `json:"requestParams" gorm:"type:text;comment:请求参数(脱敏)"`
	Result        string `json:"result" gorm:"type:varchar(50);comment:操作结果"`
	SourceIP      string `json:"sourceIp" gorm:"type:varchar(50);comment:来源IP"`
}

func (OperationAuditLog) TableName() string { return "operation_audit_logs" }

// ==================== 系统基础配置 ====================

// PricingRule 计价规则（仅配置存储，不实现计价计算）
type PricingRule struct {
	global.GVA_MODEL
	RuleName      string  `json:"ruleName" gorm:"type:varchar(100);comment:规则名称"`
	BaseFee       float64 `json:"baseFee" gorm:"type:numeric(10,2);comment:基础费用"`
	MileagePrice  float64 `json:"mileagePrice" gorm:"type:numeric(10,4);comment:里程单价(元/公里)"`
	TimePrice     float64 `json:"timePrice" gorm:"type:numeric(10,4);comment:时长单价(元/分钟)"`
	AdditionalFee float64 `json:"additionalFee" gorm:"type:numeric(10,2);comment:附加费用"`
	Status        int     `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
	Remark        string  `json:"remark" gorm:"type:varchar(255);comment:备注"`
}

func (PricingRule) TableName() string { return "pricing_rules" }

// MapVersion 地图版本管理
type MapVersion struct {
	global.GVA_MODEL
	Version   string `json:"version" gorm:"type:varchar(100);comment:版本号"`
	Publisher string `json:"publisher" gorm:"type:varchar(100);comment:发布人"`
	Status    int    `json:"status" gorm:"type:int;default:0;comment:状态 0草稿 1已发布 2已下线"`
	IsActive  bool   `json:"isActive" gorm:"type:boolean;default:false;comment:当前生效版本"`
	Remark    string `json:"remark" gorm:"type:varchar(255);comment:备注"`
}

func (MapVersion) TableName() string { return "map_versions" }

// PushTemplate 消息推送模板
type PushTemplate struct {
	global.GVA_MODEL
	TemplateName string `json:"templateName" gorm:"type:varchar(100);comment:模板名称"`
	TemplateCode string `json:"templateCode" gorm:"type:varchar(100);uniqueIndex;comment:模板编码"`
	Content      string `json:"content" gorm:"type:text;comment:模板内容(含占位符)"`
	Params       string `json:"params" gorm:"type:text;comment:参数定义(JSON)"`
	Status       int    `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
}

func (PushTemplate) TableName() string { return "push_templates" }
