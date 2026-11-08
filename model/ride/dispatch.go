package ride

// DispatchRecordListQuery 后台派单记录列表查询条件。
type DispatchRecordListQuery struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Keyword  string `json:"keyword" form:"keyword"`
	Status   string `json:"status" form:"status"`
}

// DispatchRecordItem 派单记录行，字段与 ride-service 内部接口保持一致。
type DispatchRecordItem struct {
	DispatchNo                string   `json:"dispatchNo"`
	OrderID                   int64    `json:"orderId"`
	OrderNo                   string   `json:"orderNo"`
	VehicleID                 *int64   `json:"vehicleId"`
	PlateNo                   string   `json:"plateNo"`
	DispatchStrategy          string   `json:"dispatchStrategy"`
	CandidateScore            *float64 `json:"candidateScore"`
	EstimatedPickupDistanceKM *float64 `json:"estimatedPickupDistanceKm"`
	EstimatedArrivalMin       *int32   `json:"estimatedArrivalMin"`
	Status                    int      `json:"status"`
	FailureReason             string   `json:"failureReason"`
	DispatchedAt              string   `json:"dispatchedAt"`
	RespondedAt               string   `json:"respondedAt"`
	CreatedAt                 string   `json:"createdAt"`
}

// DispatchRecordListResult 派单记录分页结果。
type DispatchRecordListResult struct {
	Items    []DispatchRecordItem `json:"items"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}
