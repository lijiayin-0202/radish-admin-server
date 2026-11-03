package trip

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ==================== 运营调度管理（新模块） ====================

// OperationDispatch 运营调度（新建表 operation_dispatch）
// 记录每一次调度派单关联的司机、车辆、订单及状态
type OperationDispatch struct {
	global.GVA_MODEL
	DispatchNo   string     `json:"dispatchNo" gorm:"type:varchar(50);comment:调度编号"`
	DriverID     uint       `json:"driverId" gorm:"comment:司机ID"`
	CarID        uint       `json:"carId" gorm:"comment:车辆ID"`
	OrderID      uint       `json:"orderId" gorm:"comment:订单ID"`
	Status       int        `json:"status" gorm:"type:int;default:0;comment:调度状态 0待调度 1调度中 2已派单 3已完成 4已取消 5异常"`
	DispatchTime *time.Time `json:"dispatchTime" gorm:"comment:调度时间"`
	Remark       string     `json:"remark" gorm:"type:varchar(500);comment:备注"`
}

func (OperationDispatch) TableName() string { return "operation_dispatch" }
