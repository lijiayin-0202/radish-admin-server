package response

// RoleWithPermissions 角色及其权限点
type RoleWithPermissions struct {
	ID            uint   `json:"id"`
	RoleName      string `json:"roleName"`
	RoleCode      string `json:"roleCode"`
	Description   string `json:"description"`
	Status        int    `json:"status"`
	PermissionIDs []uint `json:"permissionIds"`
}

// OperatorWithRoles 后台账号及其角色
type OperatorWithRoles struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	Status   int    `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

// PermissionTree 权限树节点
type PermissionTree struct {
	ID             uint             `json:"id"`
	ParentID       uint             `json:"parentId"`
	PermissionName string           `json:"permissionName"`
	PermissionType int              `json:"permissionType"`
	RoutePath      string           `json:"routePath"`
	HTTPMethod     string           `json:"httpMethod"`
	Sort           int              `json:"sort"`
	Children       []PermissionTree `json:"children"`
}

// MetricItem 指标项
type MetricItem struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// TrendPoint 趋势点
type TrendPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// DashboardOverview 数据大屏聚合指标
type DashboardOverview struct {
	// 订单指标
	OrderTotal   int64   `json:"orderTotal"`   // 下单量
	OrderDone    int64   `json:"orderDone"`    // 完成订单数
	OrderCancel  int64   `json:"orderCancel"`  // 取消订单数
	CompletionRate float64 `json:"completionRate"` // 完单率(%)
	CancelRate   float64 `json:"cancelRate"`   // 取消率(%)
	AvgOrderAmount float64 `json:"avgOrderAmount"` // 客单价

	// 车辆指标
	VehicleOnline int64 `json:"vehicleOnline"` // 在线车辆数
	VehicleTotal  int64 `json:"vehicleTotal"`  // 车辆总数

	// 客服指标
	TicketTotal   int64 `json:"ticketTotal"`   // 客服工单数
	SessionTotal  int64 `json:"sessionTotal"`  // 客服会话数

	// 审计指标
	AuditLogTotal int64 `json:"auditLogTotal"` // 审计日志数

	// 趋势
	OrderTrend   []TrendPoint `json:"orderTrend"`   // 下单量趋势
	VehicleTrend []TrendPoint `json:"vehicleTrend"` // 在线车辆趋势
}

// StationFlow 站点客流
type StationFlow struct {
	StationID   uint   `json:"stationId"`
	StationName string `json:"stationName"`
	BoardCount  int64  `json:"boardCount"`  // 上车人次
	AlightCount int64  `json:"alightCount"` // 下车人次
}

// DashboardCity 城市选项（下拉）
type DashboardCity struct {
	CityCode string `json:"cityCode"`
	CityName string `json:"cityName"`
}

// DashboardArea 区域选项（按城市联动）
type DashboardArea struct {
	ID       uint   `json:"id"`
	AreaName string `json:"areaName"`
	CityCode string `json:"cityCode"`
	CityName string `json:"cityName"`
}
