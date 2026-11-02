package bizdb

import (
	"errors"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 萝卜快跑业务数据库连接管理
// 业务事实表分布：postgres 库（乘客/订单/区域/站点/车辆/派单/指令/故障/核验/安全）
// luobo 库（客服会话/工单/FAQ/通知）；GVA 主库 gva（global.GVA_DB）

var (
	pgOnce, csOnce sync.Once
	pgDB, csDB     *gorm.DB
)

// PostgresDB 出行业务库（postgres 库）
func PostgresDB() *gorm.DB {
	pgOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(global.GVA_CONFIG.Pgsql.LinkDsn("postgres")), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err == nil {
			tunePool(db)
			pgDB = db
		}
	})
	return pgDB
}

// LuoboDB 客服业务库（luobo 库）
func LuoboDB() *gorm.DB {
	csOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(global.GVA_CONFIG.Pgsql.LinkDsn("luobo")), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err == nil {
			tunePool(db)
			csDB = db
		}
	})
	return csDB
}

// tunePool 限制连接池大小，避免大屏轮询等高频请求把 PostgreSQL 连接打满
// （曾出现 FATAL: sorry, too many clients already / SQLSTATE 53300）
func tunePool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxOpenConns(20)            // 单库最大 20 条连接
	sqlDB.SetMaxIdleConns(5)             // 最多保留 5 条空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最长存活 1 小时，避免僵死连接
}

// AutoMigrate 在 postgres 库补建缺失的业务表（仅新建表，不动已有表）
func AutoMigrate() error {
	db := PostgresDB()
	if db == nil {
		return errors.New("postgres 业务库连接不可用")
	}
	return db.AutoMigrate(
		&tripModel.Station{},
		&tripModel.DispatchRecord{},
		&tripModel.VehicleCommand{},
		&tripModel.VehicleFault{},
		&tripModel.BoardingVerification{},

		// PRD 数据字典中存在、库中缺失的表（本次补建）
		&tripModel.OrderCancellation{},
		&tripModel.OrderFeeItem{},
		&tripModel.Refund{},
		&tripModel.ServiceAreaSchedule{},
		&tripModel.VehicleStatusLog{},
		&tripModel.VehicleLocationLog{},
		&tripModel.VehicleTelemetryLog{},
		&tripModel.VehicleMediaRecord{},
	)
}

// AutoMigrateCS 在 luobo 库补建/同步客服业务表（工单处理日志等）
func AutoMigrateCS() error {
	db := LuoboDB()
	if db == nil {
		return errors.New("luobo 客服业务库连接不可用")
	}
	if err := db.AutoMigrate(
		&csModel.CsSession{},
		&csModel.CsTicket{},
		&csModel.CsFaqCategory{},
		&csModel.CsFaq{},
		&csModel.CsNotification{},
		&csModel.CsTicketLog{},
	); err != nil {
		return err
	}
	// 将工单 status 由 varchar 迁移为 int（1待处理/2处理中/3已关闭），已为 int 时自动跳过
	db.Exec("UPDATE cs_tickets SET status = CASE COALESCE(status,'') WHEN '待处理' THEN '1' WHEN '处理中' THEN '2' WHEN '已解决' THEN '3' WHEN '已关闭' THEN '3' ELSE '1' END WHERE status IS NOT NULL AND status <> ''")
	db.Exec("ALTER TABLE cs_tickets ALTER COLUMN status TYPE integer USING status::integer")
	db.Exec("ALTER TABLE cs_tickets ALTER COLUMN status SET DEFAULT 1")
	return nil
}
