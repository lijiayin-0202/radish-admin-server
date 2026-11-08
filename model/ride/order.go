package ride

// OrderListQuery 后台订单列表查询条件。
type OrderListQuery struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Status   string `json:"status" form:"status"`
	Keyword  string `json:"keyword" form:"keyword"`
}

// OrderItem 订单列表行，字段与 ride-service 内部接口保持一致。
type OrderItem struct {
	OrderID         string  `json:"orderId"`
	OrderNo         string  `json:"orderNo"`
	PassengerID     int64   `json:"passengerId"`
	Status          string  `json:"status"`
	StatusValue     int     `json:"statusValue"`
	StatusText      string  `json:"statusText"`
	OriginName      string  `json:"originName"`
	DestinationName string  `json:"destinationName"`
	PayableAmount   float64 `json:"payableAmount"`
	VehiclePlate    string  `json:"vehiclePlate"`
	CreatedAt       string  `json:"createdAt"`
}

// OrderListResult 分页结果。
type OrderListResult struct {
	Items    []OrderItem `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// DispatchRequest 后台派车请求：vehicleId 为空时按区域自动挑选可用车。
type DispatchRequest struct {
	OrderID   string `json:"orderId" binding:"required"`
	VehicleID int64  `json:"vehicleId"`
}

// CommandRequest 后台推进订单状态请求。
type CommandRequest struct {
	OrderID     string `json:"orderId" binding:"required"`
	CommandType string `json:"commandType" binding:"required"`
}

// DispatchResult 派车结果。
type DispatchResult struct {
	Order   OrderItem `json:"order"`
	Vehicle struct {
		VehicleID int64  `json:"vehicleId"`
		PlateNo   string `json:"plateNo"`
		Model     string `json:"model"`
	} `json:"vehicle"`
}
