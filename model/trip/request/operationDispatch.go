package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// OperationDispatchSearch 运营调度查询（分页 + 筛选）
type OperationDispatchSearch struct {
	request.PageInfo
	DispatchNo string `json:"dispatchNo" form:"dispatchNo"` // 调度编号（模糊匹配）
	DriverID   uint   `json:"driverId" form:"driverId"`     // 司机ID（0=不筛选）
	CarID      uint   `json:"carId" form:"carId"`           // 车辆ID（0=不筛选）
	OrderID    uint   `json:"orderId" form:"orderId"`       // 订单ID（0=不筛选）
	Status     *int   `json:"status" form:"status"`         // 调度状态（nil=不筛选，0待调度 是合法值）
	StartTime  string `json:"startTime" form:"startTime"`   // 调度时间起（含）
	EndTime    string `json:"endTime" form:"endTime"`       // 调度时间止（含）
}
