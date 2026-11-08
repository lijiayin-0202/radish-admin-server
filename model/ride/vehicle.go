package ride

// VehicleListQuery 后台车辆列表查询条件。
type VehicleListQuery struct {
	Page            int    `json:"page" form:"page"`
	PageSize        int    `json:"pageSize" form:"pageSize"`
	Keyword         string `json:"keyword" form:"keyword"`
	OperationStatus string `json:"operationStatus" form:"operationStatus"`
	OnlineStatus    string `json:"onlineStatus" form:"onlineStatus"`
}

// VehicleItem 车辆列表行，字段与 ride-service 内部接口保持一致。
type VehicleItem struct {
	VehicleID       int64    `json:"vehicleId"`
	VehicleNo       string   `json:"vehicleNo"`
	VIN             string   `json:"vin"`
	PlateNo         string   `json:"plateNo"`
	Model           string   `json:"model"`
	SeatCount       int      `json:"seatCount"`
	ServiceAreaID   *int64   `json:"serviceAreaId"`
	ServiceAreaName string   `json:"serviceAreaName"`
	OperationStatus int      `json:"operationStatus"`
	OnlineStatus    int      `json:"onlineStatus"`
	BatteryPercent  *float64 `json:"batteryPercent"`
	SoftwareVersion string   `json:"softwareVersion"`
	LastReportedAt  string   `json:"lastReportedAt"`
	CreatedAt       string   `json:"createdAt"`
}

// VehicleListResult 车辆分页结果。
type VehicleListResult struct {
	Items    []VehicleItem `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// VehicleDetailResult 车辆详情结果。
type VehicleDetailResult struct {
	Vehicle VehicleItem `json:"vehicle"`
}
