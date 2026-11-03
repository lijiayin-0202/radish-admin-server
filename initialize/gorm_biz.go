package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformModel "github.com/flipped-aurora/gin-vue-admin/server/model/platform"
	"github.com/flipped-aurora/gin-vue-admin/server/service/bizdb"
	"go.uber.org/zap"
)

func bizModel() error {
	db := global.GVA_DB
	// 平台配置、权限与统计大屏 业务表
	err := db.AutoMigrate(
		platformModel.Agreement{},
		platformModel.Announcement{},
		platformModel.PlatformSetting{},
		platformModel.Role{},
		platformModel.Permission{},
		platformModel.RolePermission{},
		platformModel.OperatorUser{},
		platformModel.OperatorUserRole{},
		platformModel.OperationAuditLog{},
		platformModel.PricingRule{},
		platformModel.MapVersion{},
		platformModel.PushTemplate{},
	)
	if err != nil {
		return err
	}
	// 出行业务库（postgres 库）补建缺失业务表：站点/派单/指令/故障/核验/取消/费用/退款/车辆日志等
	if err := bizdb.AutoMigrate(); err != nil {
		global.GVA_LOG.Error("trip biz tables migrate failed", zap.Error(err))
	}
	// 客服业务库（luobo 库）补建/同步客服业务表：会话/工单/FAQ/通知/工单日志
	if err := bizdb.AutoMigrateCS(); err != nil {
		global.GVA_LOG.Error("cs biz tables migrate failed", zap.Error(err))
	}
	return nil
}
