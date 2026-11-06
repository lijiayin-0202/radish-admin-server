package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// LuoboPgSql 初始化第二个 PostgreSQL 连接 luoboPgsql，连接成功后赋值给 global.LUOBO_DB
// 说明：
//  1. 只建立连接与连接池配置，不做任何 AutoMigrate，表结构一律由外部 SQL 脚本维护
//  2. 不影响 global.GVA_DB 及系统库相关初始化流程
//  3. 未配置或连接不可用时只记录日志，不阻断主库启动
func LuoboPgSql() {
	p := global.GVA_CONFIG.LuoboPgsql
	if p.Dbname == "" || p.Path == "" {
		global.GVA_LOG.Warn("luoboPgsql 未配置或配置不完整（db-name / path 为空），跳过第二个 PostgreSQL 连接初始化")
		return
	}

	db := GormPgSqlByConfig(p)
	if db == nil {
		global.GVA_LOG.Warn("luoboPgsql 初始化返回空连接，global.LUOBO_DB 保持为 nil")
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.GVA_LOG.Error("luoboPgsql 获取底层 *sql.DB 失败", zap.Error(err))
		return
	}
	if err = sqlDB.Ping(); err != nil {
		global.GVA_LOG.Error("luoboPgsql 连通性检查失败", zap.Error(err))
		return
	}

	global.LUOBO_DB = db
	global.GVA_LOG.Info("luoboPgsql 初始化成功",
		zap.String("path", p.Path),
		zap.String("port", p.Port),
		zap.String("dbname", p.Dbname))
}
