package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// TripSearch 出行业务通用查询（分页 + 关键词 + 可选状态 + 城市/区域/时间维度）
type TripSearch struct {
	request.PageInfo
	Keyword   string `json:"keyword" form:"keyword"`
	Status    *int   `json:"status" form:"status"`
	City      string `json:"city" form:"city"`           // 城市编码（空=全部城市）
	Area      uint   `json:"area" form:"area"`           // 运营区域ID（空=该城市全部区域）
	StartTime string `json:"startTime" form:"startTime"` // 时间范围起（含），格式 YYYY-MM-DD HH:MM:SS
	EndTime   string `json:"endTime" form:"endTime"`     // 时间范围止（含）
}

// OrderSearch 订单查询
type OrderSearch struct {
	request.PageInfo
	Keyword     string `json:"keyword" form:"keyword"`
	OrderStatus *int   `json:"orderStatus" form:"orderStatus"`
	City        string `json:"city" form:"city"`
	Area        uint   `json:"area" form:"area"`
	StartTime   string `json:"startTime" form:"startTime"`
	EndTime     string `json:"endTime" form:"endTime"`
}

// OrderLogQuery 订单状态日志查询
type OrderLogQuery struct {
	OrderID uint `json:"orderId" form:"orderId"`
}
