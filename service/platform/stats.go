package platform

import (
	"errors"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformReq "github.com/flipped-aurora/gin-vue-admin/server/model/platform/request"
	platformRes "github.com/flipped-aurora/gin-vue-admin/server/model/platform/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/bizdb"
	"gorm.io/gorm"
)

// openDB 获取业务库连接。
// 复用 bizdb 中 sync.Once 缓存的单连接池（postgres/luobo），
// 避免每次请求都 gorm.Open 新建池，导致大屏轮询把 PostgreSQL 连接打满
// （曾报 FATAL: sorry, too many clients already / SQLSTATE 53300）。
func openDB(dbname string) (*gorm.DB, error) {
	var db *gorm.DB
	switch dbname {
	case "postgres":
		db = bizdb.PostgresDB()
	case "luobo":
		db = bizdb.LuoboDB()
	default:
		return nil, errors.New("未知数据库: " + dbname)
	}
	if db == nil {
		return nil, errors.New(dbname + " 业务库连接不可用")
	}
	return db, nil
}

// buildOrderScope 根据筛选条件构造 trip_orders 的 WHERE 片段与参数
// 支持：城市（city_code 关联 service_areas）/ 区域（service_area_id）/ 时间范围
func buildOrderScope(q platformReq.DashboardQuery) (string, []interface{}) {
	where := "1=1"
	args := make([]interface{}, 0)
	// 区域优先：直接按 service_area_id 过滤；否则按城市过滤其下所有区域
	if q.Area > 0 {
		where += " AND service_area_id = ?"
		args = append(args, q.Area)
	} else if q.City != "" {
		where += " AND service_area_id IN (SELECT id FROM service_areas WHERE city_code = ?)"
		args = append(args, q.City)
	}
	if q.StartTime != "" && q.EndTime != "" {
		where += " AND created_at >= ? AND created_at <= ?"
		args = append(args, q.StartTime, q.EndTime)
	}
	return where, args
}

// buildVehicleScope 车辆筛选片段（与订单同维度：城市/区域）
func buildVehicleScope(q platformReq.DashboardQuery) (string, []interface{}) {
	where := "1=1"
	args := make([]interface{}, 0)
	if q.Area > 0 {
		where += " AND service_area_id = ?"
		args = append(args, q.Area)
	} else if q.City != "" {
		where += " AND service_area_id IN (SELECT id FROM service_areas WHERE city_code = ?)"
		args = append(args, q.City)
	}
	return where, args
}

// GetDashboardOverview 数据大屏聚合指标
func (s *PlatformService) GetDashboardOverview(q platformReq.DashboardQuery) (res platformRes.DashboardOverview, err error) {
	// 订单数据在 postgres 库
	ordersDB, err := openDB("postgres")
	if err != nil {
		return res, err
	}
	vehiclesDB, err := openDB("postgres")
	if err != nil {
		return res, err
	}
	// 客服数据在 luobo 库
	csDB, err := openDB("luobo")
	if err != nil {
		return res, err
	}

	orderWhere, args := buildOrderScope(q)
	vehWhere, vargs := buildVehicleScope(q)

	// 订单指标
	ordersDB.Raw("SELECT count(*) FROM trip_orders WHERE "+orderWhere, args...).Scan(&res.OrderTotal)
	ordersDB.Raw("SELECT count(*) FROM trip_orders WHERE order_status = 5 AND "+orderWhere, args...).Scan(&res.OrderDone)
	ordersDB.Raw("SELECT count(*) FROM trip_orders WHERE order_status = 6 AND "+orderWhere, args...).Scan(&res.OrderCancel)
	ordersDB.Raw("SELECT COALESCE(AVG(payable_amount),0) FROM trip_orders WHERE order_status = 5 AND "+orderWhere, args...).Scan(&res.AvgOrderAmount)

	if res.OrderTotal > 0 {
		res.CompletionRate = float64(res.OrderDone) / float64(res.OrderTotal) * 100
		res.CancelRate = float64(res.OrderCancel) / float64(res.OrderTotal) * 100
	}

	// 车辆指标（status = 1 视为在线），按城市/区域过滤
	vehiclesDB.Raw("SELECT count(*) FROM vehicles WHERE "+vehWhere, vargs...).Scan(&res.VehicleTotal)
	vehiclesDB.Raw("SELECT count(*) FROM vehicles WHERE status = 1 AND "+vehWhere, vargs...).Scan(&res.VehicleOnline)

	// 客服指标（全量，跨城市不分城市维度）
	csDB.Raw("SELECT count(*) FROM cs_tickets").Scan(&res.TicketTotal)
	csDB.Raw("SELECT count(*) FROM cs_sessions").Scan(&res.SessionTotal)

	// 审计指标（gva 主库，全量）
	global.GVA_DB.Raw("SELECT count(*) FROM operation_audit_logs").Scan(&res.AuditLogTotal)

	// 下单量趋势：同一自然日内按小时聚合（补零成 24 个点），跨天按天聚合。
	// 这样一次刷新即返回多个点，且各时段数量随真实订单变化。
	type trendRow struct {
		Date  string
		Count int64
	}
	sameDay := len(q.StartTime) >= 10 && len(q.EndTime) >= 10 && q.StartTime[:10] == q.EndTime[:10]
	groupFmt := "YYYY-MM-DD"
	if sameDay {
		groupFmt = "HH24"
	}
	var orderTrend []trendRow
	ordersDB.Raw("SELECT to_char(created_at,'"+groupFmt+"') AS date, count(*) AS count FROM trip_orders WHERE "+orderWhere+" GROUP BY date ORDER BY date", args...).Scan(&orderTrend)
	if sameDay {
		// 补零：输出 00:00 ~ 23:00 共 24 个点，缺失小时记 0
		bucket := make(map[string]int64, len(orderTrend))
		for _, r := range orderTrend {
			bucket[r.Date] = r.Count
		}
		for h := 0; h < 24; h++ {
			key := fmt.Sprintf("%02d", h)
			res.OrderTrend = append(res.OrderTrend, platformRes.TrendPoint{Date: key + ":00", Value: float64(bucket[key])})
		}
	} else {
		for _, r := range orderTrend {
			res.OrderTrend = append(res.OrderTrend, platformRes.TrendPoint{Date: r.Date, Value: float64(r.Count)})
		}
	}

	return res, nil
}

// GetCities 城市下拉选项（从运营区域表去重得到）
func (s *PlatformService) GetCities() ([]platformRes.DashboardCity, error) {
	db, err := openDB("postgres")
	if err != nil {
		return nil, err
	}
	var list []platformRes.DashboardCity
	err = db.Raw("SELECT city_code AS city_code, MAX(city_name) AS city_name FROM service_areas WHERE city_code <> '' GROUP BY city_code ORDER BY city_code").Scan(&list).Error
	return list, err
}

// GetAreas 区域下拉选项（按城市联动；city 为空则返回全部）
func (s *PlatformService) GetAreas(cityCode string) ([]platformRes.DashboardArea, error) {
	db, err := openDB("postgres")
	if err != nil {
		return nil, err
	}
	q := "SELECT id, area_name, city_code, city_name FROM service_areas WHERE status = 1"
	args := make([]interface{}, 0)
	if cityCode != "" {
		q += " AND city_code = ?"
		args = append(args, cityCode)
	}
	q += " ORDER BY id"
	var list []platformRes.DashboardArea
	err = db.Raw(q, args...).Scan(&list).Error
	return list, err
}

// VehiclePoint 车辆位置点
type VehiclePoint struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Status    int     `json:"status"`
}

// GetVehicleHeatmap 车辆分布热力图数据
func (s *PlatformService) GetVehicleHeatmap() ([]VehiclePoint, error) {
	db, err := openDB("postgres")
	if err != nil {
		return nil, err
	}
	var points []VehiclePoint
	err = db.Raw("SELECT current_longitude AS longitude, current_latitude AS latitude, status FROM vehicles WHERE current_longitude IS NOT NULL AND current_latitude IS NOT NULL").Scan(&points).Error
	return points, err
}

// GetStationFlow 站点上下客人次（按订单上下车点名称聚合，携带城市/区域筛选）
func (s *PlatformService) GetStationFlow(q platformReq.DashboardQuery) ([]platformRes.StationFlow, error) {
	db, err := openDB("postgres")
	if err != nil {
		return nil, err
	}
	orderWhere, args := buildOrderScope(q)
	// UNION ALL 的两段 SQL 各使用一次 orderWhere，占位符参数需拼两份，
	// 否则第二段占位符无对应参数，PostgreSQL 会报 syntax error at or near "AND"。
	unionArgs := append(append([]interface{}{}, args...), args...)
	type row struct {
		Name        string
		BoardCount  int64
		AlightCount int64
	}
	var rows []row
	err = db.Raw(`SELECT COALESCE(pickup_name,'未知') AS name, count(*) AS board_count, 0 AS alight_count
		FROM trip_orders WHERE `+orderWhere+` GROUP BY pickup_name
		UNION ALL
		SELECT COALESCE(destination_name,'未知') AS name, 0 AS board_count, count(*) AS alight_count
		FROM trip_orders WHERE `+orderWhere+` GROUP BY destination_name`, unionArgs...).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	// 合并同名站点
	merged := map[string]*platformRes.StationFlow{}
	for _, r := range rows {
		m, ok := merged[r.Name]
		if !ok {
			m = &platformRes.StationFlow{StationName: r.Name}
			merged[r.Name] = m
		}
		m.BoardCount += r.BoardCount
		m.AlightCount += r.AlightCount
	}
	result := make([]platformRes.StationFlow, 0, len(merged))
	for _, m := range merged {
		result = append(result, *m)
	}
	return result, nil
}
