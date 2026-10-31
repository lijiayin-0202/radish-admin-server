package platform

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PlatformRouter struct{}

// InitPlatformRouter 注册平台配置、权限与统计大屏路由
// 路由前缀：/platform
func (s *PlatformRouter) InitPlatformRouter(Router *gin.RouterGroup) {
	router := Router.Group("platform").Use(middleware.OperationRecord())
	readonly := Router.Group("platform")
	api := v1.ApiGroupApp.PlatformApiGroup.PlatformApi
	{
		// 协议
		router.POST("agreement", api.CreateAgreement)
		router.PUT("agreement", api.UpdateAgreement)
		router.DELETE("agreement/:id", api.DeleteAgreement)
		router.DELETE("agreements", api.DeleteAgreementsByIds)

		// 公告
		router.POST("announcement", api.CreateAnnouncement)
		router.PUT("announcement", api.UpdateAnnouncement)
		router.DELETE("announcement/:id", api.DeleteAnnouncement)
		router.DELETE("announcements", api.DeleteAnnouncementsByIds)

		// 平台参数
		router.POST("setting", api.CreateSetting)
		router.PUT("setting", api.UpdateSetting)
		router.DELETE("setting/:id", api.DeleteSetting)
		router.DELETE("settings", api.DeleteSettingsByIds)

		// 角色
		router.POST("role", api.CreateRole)
		router.PUT("role", api.UpdateRole)
		router.DELETE("role/:id", api.DeleteRole)
		router.POST("role/permissions", api.AssignRolePermissions)

		// 权限点
		router.POST("permission", api.CreatePermission)
		router.PUT("permission", api.UpdatePermission)
		router.DELETE("permission/:id", api.DeletePermission)

		// 后台账号
		router.POST("operator", api.CreateOperator)
		router.PUT("operator", api.UpdateOperator)
		router.DELETE("operator/:id", api.DeleteOperator)
		router.POST("operator/roles", api.AssignOperatorRoles)

		// 计价规则
		router.POST("pricingRule", api.CreatePricingRule)
		router.PUT("pricingRule", api.UpdatePricingRule)
		router.DELETE("pricingRule/:id", api.DeletePricingRule)
		router.DELETE("pricingRules", api.DeletePricingRulesByIds)

		// 地图版本
		router.POST("mapVersion", api.CreateMapVersion)
		router.PUT("mapVersion", api.UpdateMapVersion)
		router.DELETE("mapVersion/:id", api.DeleteMapVersion)
		router.POST("mapVersion/:id/activate", api.ActivateMapVersion)

		// 消息推送模板
		router.POST("pushTemplate", api.CreatePushTemplate)
		router.PUT("pushTemplate", api.UpdatePushTemplate)
		router.DELETE("pushTemplate/:id", api.DeletePushTemplate)
		router.DELETE("pushTemplates", api.DeletePushTemplatesByIds)
	}
	{
		// 协议
		readonly.GET("agreement/list", api.GetAgreementList)
		readonly.GET("agreement/:id", api.GetAgreementById)

		// 公告
		readonly.GET("announcement/list", api.GetAnnouncementList)

		// 平台参数
		readonly.GET("setting/list", api.GetSettingList)

		// 角色
		readonly.GET("role/list", api.GetRoleList)
		readonly.GET("role/:id/permissions", api.GetRolePermissions)

		// 权限点
		readonly.GET("permission/list", api.GetPermissionList)
		readonly.GET("permission/tree", api.GetPermissionTree)

		// 后台账号
		readonly.GET("operator/list", api.GetOperatorList)
		readonly.GET("operator/:id/roles", api.GetOperatorRoles)

		// 审计日志
		readonly.GET("auditLog/list", api.GetAuditLogList)

		// 计价规则
		readonly.GET("pricingRule/list", api.GetPricingRuleList)

		// 地图版本
		readonly.GET("mapVersion/list", api.GetMapVersionList)

		// 消息模板
		readonly.GET("pushTemplate/list", api.GetPushTemplateList)

		// 统计大屏
		readonly.GET("dashboard/overview", api.GetDashboardOverview)
		readonly.GET("dashboard/vehicle-heatmap", api.GetVehicleHeatmap)
		readonly.GET("dashboard/station-flow", api.GetStationFlow)
		readonly.GET("dashboard/cities", api.GetDashboardCities)
		readonly.GET("dashboard/areas", api.GetDashboardAreas)
	}
}
