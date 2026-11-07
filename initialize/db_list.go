package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const sys = "system"

func DBList() {
	dbMap := make(map[string]*gorm.DB)
	for _, info := range global.GVA_CONFIG.DBList {
		if info.Disable {
			continue
		}
		switch info.Type {
		case "mysql":
			dbMap[info.AliasName] = GormMysqlByConfig(config.Mysql{GeneralDB: info.GeneralDB})
		case "mssql":
			dbMap[info.AliasName] = GormMssqlByConfig(config.Mssql{GeneralDB: info.GeneralDB})
		case "pgsql":
			dbMap[info.AliasName] = GormPgSqlByConfig(config.Pgsql{GeneralDB: info.GeneralDB})
		case "oracle":
			dbMap[info.AliasName] = GormOracleByConfig(config.Oracle{GeneralDB: info.GeneralDB})
		default:
			continue
		}
	}
	// 做特殊判断,是否有迁移
	// 适配低版本迁移多数据库版本
	if sysDB, ok := dbMap[sys]; ok {
		global.GVA_DB = sysDB
	}
	global.GVA_DBList = dbMap
}

// UseGvaPgAsMainDB 把 GVA 主句柄指向 db-list 中的 gvaPg 实例。
// 注入后，框架全部权限表（sys_*、casbin_rule）、启动建表、基础数据初始化
// 以及业务表的读写都统一走该 PG 实例；原有 MySQL 代码与配置保持不变。
// 必须在 DBList() 之后、RegisterTables() 之前调用。
func UseGvaPgAsMainDB() {
	if global.GVA_DBList == nil || global.GVA_DBList[global.GvaPgDB] == nil {
		global.GVA_LOG.Warn("db-list 中未找到可用的 gvaPg，保持原有主库连接")
		return
	}
	global.GVA_DB = global.MustGetGlobalDBByDBName(global.GvaPgDB)
	for _, info := range global.GVA_CONFIG.DBList {
		if info.AliasName == global.GvaPgDB && !info.Disable {
			dbname := info.Dbname
			global.GVA_ACTIVE_DBNAME = &dbname
			break
		}
	}
	global.GVA_LOG.Info("GVA 系统库与业务表统一使用 db-list[gvaPg]",
		zap.String("dbname", func() string {
			if global.GVA_ACTIVE_DBNAME != nil {
				return *global.GVA_ACTIVE_DBNAME
			}
			return ""
		}()))
}
