package bizdb

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	"gorm.io/gorm"
)

// 萝卜快跑业务表统一落在 db-list[gvaPg] 指向的 gva 库中：
// 运营调度/乘客/订单等由 PostgresDB() 承载，客服会话/工单/FAQ 由 LuoboDB() 承载。
// 两个函数返回同一个连接（global.MustGetGlobalDBByDBName(global.GvaPgDB)），
// 不再分别连接 postgres / luobo 独立库；连接池由 DBList() 统一管理。
// 注：PostgresDB / LuoboDB 仅为兼容既有 service 调用而保留的语义别名。

// PostgresDB 运营调度、乘客相关业务表
func PostgresDB() *gorm.DB {
	return global.MustGetGlobalDBByDBName(global.GvaPgDB)
}

// LuoboDB 客服业务表（会话/工单/FAQ/通知）
func LuoboDB() *gorm.DB {
	return global.MustGetGlobalDBByDBName(global.GvaPgDB)
}

// AutoMigrate 在 db-list[gvaPg] 指向的 gva 库补建缺失的业务表
func AutoMigrate() error {
	db := PostgresDB()
	if db == nil {
		return errors.New("gvaPg 业务库连接不可用")
	}
	return db.AutoMigrate(tripBusinessModels()...)
}

func tripBusinessModels() []interface{} {
	return []interface{}{
		// 核心业务表：原先由外部 SQL 建在 postgres 库，统一到 gva 库后必须纳入迁移，
		// 否则订单/乘客/区域/车辆等页面会报 relation does not exist、列表无数据，
		// passengers 也会因缺 deleted_at 报 SQLSTATE 42703。
		&tripModel.Passenger{},
		&tripModel.TripOrder{},
		&tripModel.TripStatusLog{},
		&tripModel.ServiceArea{},
		&tripModel.Vehicle{},
		&tripModel.SafetyIncident{},

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
		&tripModel.OperationDispatch{},
	}
}

// AutoMigrateCS 在 db-list[gvaPg] 指向的 gva 库补建/同步客服业务表
func AutoMigrateCS() error {
	db := LuoboDB()
	if db == nil {
		return errors.New("gvaPg 客服业务库连接不可用")
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
	// 历史外部 SQL 建的 cs_tickets.status 为 varchar（中文状态），而 AutoMigrate
	// 新建的表 status 已是 int。仅当列确为文本类型时才做转换，避免每次启动
	// 在 int 列上执行 COALESCE(status,'') 产生 PG 类型错误日志。
	if !csTicketStatusIsText(db) {
		return nil
	}
	db.Exec("UPDATE cs_tickets SET status = CASE COALESCE(status,'') WHEN '待处理' THEN '1' WHEN '处理中' THEN '2' WHEN '已解决' THEN '3' WHEN '已关闭' THEN '3' ELSE '1' END WHERE status IS NOT NULL AND status <> ''")
	db.Exec("ALTER TABLE cs_tickets ALTER COLUMN status TYPE integer USING status::integer")
	db.Exec("ALTER TABLE cs_tickets ALTER COLUMN status SET DEFAULT 1")
	return nil
}

func csTicketStatusIsText(db *gorm.DB) bool {
	types, err := db.Migrator().ColumnTypes(&csModel.CsTicket{})
	if err != nil {
		return false
	}
	for _, column := range types {
		if column.Name() != "status" {
			continue
		}
		t := strings.ToLower(column.DatabaseTypeName())
		return strings.Contains(t, "char") || strings.Contains(t, "text")
	}
	return false
}
