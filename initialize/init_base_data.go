package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

// InitBaseData 启动时幂等初始化原生 GVA 基础数据。
// 复用 source/system 注册的各 initializer：sys_apis、sys_ignore_apis、
// casbin_rule、sys_authorities、sys_base_menus（40 条原生菜单）、sys_users、
// 角色-菜单授权、字典等，排序/路由/图标/层级均与原生 GVA 一致。
//
// 前提：gva 数据库已经存在（Navicat 预建），且 UseGvaPgAsMainDB() 已执行。
// 这里只做表迁移和数据灌入，不创建数据库、不回写配置文件。
// 各 initializer 的 DataInserted() 会跳过已初始化的数据，可重复启动。
func InitBaseData() {
	ctx := context.Background()
	if err := (&systemService.InitDBService{}).InitTablesAndData(ctx); err != nil {
		global.GVA_LOG.Error("init base data failed", zap.Error(err))
		return
	}
	global.GVA_LOG.Info("init base data finished on db-list[gvaPg]")
}
