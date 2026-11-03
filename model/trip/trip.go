package trip

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// BizModel 业务库基础模型（对应无 deleted_at 的表）
type BizModel struct {
	ID        uint      `gorm:"primarykey" json:"ID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ==================== 乘客（模块1 后台视角） ====================

// Passenger 乘客账号
type Passenger struct {
	global.GVA_MODEL
	UserID    uint   `json:"userId" gorm:"comment:关联用户ID"`
	Mobile    string `json:"mobile" gorm:"type:varchar(20);comment:手机号"`
	Nickname  string `json:"nickname" gorm:"type:varchar(100);comment:昵称"`
	AvatarURL string `json:"avatarUrl" gorm:"type:varchar(255);comment:头像"`
	Status    int    `json:"status" gorm:"type:int;default:1;comment:状态 1正常 0冻结"`
}

func (Passenger) TableName() string { return "passengers" }

// ==================== 订单（模块2 后台视角） ====================

// TripOrder 出行订单
type TripOrder struct {
	BizModel
	OrderNo           string     `json:"orderNo" gorm:"type:varchar(50);comment:订单号"`
	PassengerID       uint       `json:"passengerId" gorm:"comment:乘客ID"`
	VehicleID         uint       `json:"vehicleId" gorm:"comment:车辆ID"`
	ServiceAreaID     uint       `json:"serviceAreaId" gorm:"comment:运营区域ID"`
	OrderStatus       int        `json:"orderStatus" gorm:"comment:状态 0待派车 1派车中 2接驾中 3待上车 4行程中 5已完成 6已取消 7异常关闭"`
	PickupName         string     `json:"pickupName" gorm:"type:varchar(200);comment:上车点"`
	PickupAddress      string     `json:"pickupAddress" gorm:"type:varchar(300);comment:上车点地址"`
	PickupLongitude    float64    `json:"pickupLongitude" gorm:"type:numeric(10,6);comment:上车点经度"`
	PickupLatitude     float64    `json:"pickupLatitude" gorm:"type:numeric(10,6);comment:上车点纬度"`
	DestinationName    string     `json:"destinationName" gorm:"type:varchar(200);comment:下车点"`
	DestinationAddress string     `json:"destinationAddress" gorm:"type:varchar(300);comment:下车点地址"`
	DestinationLongitude float64  `json:"destinationLongitude" gorm:"type:numeric(10,6);comment:下车点经度"`
	DestinationLatitude  float64  `json:"destinationLatitude" gorm:"type:numeric(10,6);comment:下车点纬度"`
	EstimatedDistanceM int        `json:"estimatedDistanceM" gorm:"type:int;comment:预估里程(米)"`
	EstimatedDurationS int        `json:"estimatedDurationS" gorm:"type:int;comment:预估时长(秒)"`
	ActualDistanceM    int        `json:"actualDistanceM" gorm:"type:int;comment:实际里程(米)"`
	ActualDurationS    int        `json:"actualDurationS" gorm:"type:int;comment:实际时长(秒)"`
	EstimatedAmount    float64    `json:"estimatedAmount" gorm:"type:numeric(10,2);comment:预估金额"`
	OriginalAmount     float64    `json:"originalAmount" gorm:"type:numeric(10,2);comment:原价"`
	DiscountAmount     float64    `json:"discountAmount" gorm:"type:numeric(10,2);comment:优惠金额"`
	PayableAmount      float64    `json:"payableAmount" gorm:"type:numeric(10,2);comment:应付金额"`
	PassengerCouponID  *uint      `json:"passengerCouponId" gorm:"comment:乘客优惠券ID"`
	RequestedAt        *time.Time `json:"requestedAt" gorm:"comment:请求时间"`
	AssignedAt         *time.Time `json:"assignedAt" gorm:"comment:派单时间"`
	ArrivedAt          *time.Time `json:"arrivedAt" gorm:"comment:车辆到达时间"`
	StartedAt          *time.Time `json:"startedAt" gorm:"comment:行程开始"`
	CompletedAt        *time.Time `json:"completedAt" gorm:"comment:行程完成"`
	CancelledAt        *time.Time `json:"cancelledAt" gorm:"comment:取消时间"`
	CancelReason       string     `json:"cancelReason" gorm:"type:varchar(255);comment:取消原因"`
}

func (TripOrder) TableName() string { return "trip_orders" }

// TripStatusLog 订单状态流转日志（只读）
type TripStatusLog struct {
	ID           uint      `gorm:"primarykey" json:"ID"`
	OrderID      uint      `json:"orderId"`
	FromStatus   int       `json:"fromStatus"`
	ToStatus     int       `json:"toStatus"`
	OperatorType string    `json:"operatorType" gorm:"type:varchar(50)"`
	OperatorID   uint      `json:"operatorId"`
	Remark       string    `json:"remark" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (TripStatusLog) TableName() string { return "trip_status_logs" }

// ==================== 运营区域 / 站点（模块3） ====================

// ServiceArea 运营区域
type ServiceArea struct {
	BizModel
	CityCode      string `json:"cityCode" gorm:"type:varchar(20);comment:城市编码"`
	CityName      string `json:"cityName" gorm:"type:varchar(50);comment:城市名称"`
	AreaName      string `json:"areaName" gorm:"type:varchar(100);comment:区域名称"`
	PolygonGeojson string `json:"polygonGeojson" gorm:"type:jsonb;comment:区域边界GeoJSON"`
	Status        int    `json:"status" gorm:"type:smallint;default:1;comment:状态 1运营 0停运"`
}

func (ServiceArea) TableName() string { return "service_areas" }

// Station 站点（新建表，建在 postgres 库）
type Station struct {
	global.GVA_MODEL
	StationNo    string  `json:"stationNo" gorm:"type:varchar(50);comment:站点编号"`
	StationName  string  `json:"stationName" gorm:"type:varchar(100);comment:站点名称"`
	Address      string  `json:"address" gorm:"type:varchar(300);comment:站点地址"`
	Longitude    float64 `json:"longitude" gorm:"type:numeric(10,6);comment:经度"`
	Latitude     float64 `json:"latitude" gorm:"type:numeric(10,6);comment:纬度"`
	StationType  string  `json:"stationType" gorm:"type:varchar(50);comment:站点类型"`
	BoardGuide   string  `json:"boardGuide" gorm:"type:varchar(300);comment:上车指引"`
	AlightGuide  string  `json:"alightGuide" gorm:"type:varchar(300);comment:下车指引"`
	ServiceAreaID uint   `json:"serviceAreaId" gorm:"comment:所属运营区域"`
	Status       int     `json:"status" gorm:"type:int;default:1;comment:状态 1启用 0停用"`
}

func (Station) TableName() string { return "stations" }

// ==================== 车辆 / 派单（模块3） ====================

// Vehicle 车辆
type Vehicle struct {
	global.GVA_MODEL
	VehicleNo       string     `json:"vehicleNo" gorm:"type:varchar(50);comment:车辆编号"`
	PlateNo         string     `json:"plateNo" gorm:"type:varchar(20);comment:车牌号"`
	Brand           string     `json:"brand" gorm:"type:varchar(50);comment:品牌"`
	Model           string     `json:"model" gorm:"type:varchar(50);comment:型号"`
	SeatCount       int        `json:"seatCount" gorm:"comment:座位数"`
	Color           string     `json:"color" gorm:"type:varchar(20);comment:颜色"`
	ServiceAreaID   uint       `json:"serviceAreaId" gorm:"comment:所属运营区域"`
	Status          int        `json:"status" gorm:"type:int;default:1;comment:状态 1在线运营 0离线"`
	CurrentLongitude float64   `json:"currentLongitude" gorm:"type:numeric(10,6);comment:当前经度"`
	CurrentLatitude  float64   `json:"currentLatitude" gorm:"type:numeric(10,6);comment:当前纬度"`
	LastLocationAt  *time.Time `json:"lastLocationAt" gorm:"comment:最后位置上报时间"`
}

func (Vehicle) TableName() string { return "vehicles" }

// DispatchRecord 派单记录（新建表）
type DispatchRecord struct {
	global.GVA_MODEL
	DispatchNo      string  `json:"dispatchNo" gorm:"type:varchar(50);comment:派单号"`
	OrderID         uint    `json:"orderId" gorm:"comment:订单ID"`
	VehicleID       uint    `json:"vehicleId" gorm:"comment:车辆ID"`
	DispatchStrategy string `json:"dispatchStrategy" gorm:"type:varchar(100);comment:派单策略"`
	CandidateScore  float64 `json:"candidateScore" gorm:"type:numeric(10,2);comment:候选评分"`
	EtaDistanceM    int64   `json:"etaDistanceM" gorm:"comment:预计接驾距离(米)"`
	EtaDurationS    int64   `json:"etaDurationS" gorm:"comment:预计接驾时长(秒)"`
	DispatchStatus  string  `json:"dispatchStatus" gorm:"type:varchar(50);comment:状态 已创建/已下发/车辆确认/失败/超时"`
	FailReason      string  `json:"failReason" gorm:"type:varchar(255);comment:失败原因"`
}

func (DispatchRecord) TableName() string { return "dispatch_records" }

// ==================== 车辆端（模块4） ====================

// VehicleCommand 车辆指令（新建表）
type VehicleCommand struct {
	global.GVA_MODEL
	CommandNo     string `json:"commandNo" gorm:"type:varchar(50);comment:指令编号"`
	VehicleID     uint   `json:"vehicleId" gorm:"comment:车辆ID"`
	OrderID       uint   `json:"orderId" gorm:"comment:关联订单ID"`
	CommandType   string `json:"commandType" gorm:"type:varchar(50);comment:指令类型 接驾/开始行程/结束行程/靠边停车/远程协助/配置同步"`
	CommandParams string `json:"commandParams" gorm:"type:text;comment:指令参数(JSON)"`
	CommandStatus string `json:"commandStatus" gorm:"type:varchar(50);comment:状态 已发送/已确认/成功/失败"`
	ResultMessage string `json:"resultMessage" gorm:"type:varchar(255);comment:执行结果"`
	OperatorID    uint   `json:"operatorId" gorm:"comment:下发操作人"`
}

func (VehicleCommand) TableName() string { return "vehicle_commands" }

// VehicleFault 车辆故障（新建表）
type VehicleFault struct {
	global.GVA_MODEL
	FaultNo     string     `json:"faultNo" gorm:"type:varchar(50);comment:故障编号"`
	VehicleID   uint       `json:"vehicleId" gorm:"comment:车辆ID"`
	FaultLevel  string     `json:"faultLevel" gorm:"type:varchar(20);comment:故障等级"`
	FaultType   string     `json:"faultType" gorm:"type:varchar(100);comment:故障类型"`
	Description string     `json:"description" gorm:"type:text;comment:故障描述"`
	FaultStatus string     `json:"faultStatus" gorm:"type:varchar(50);comment:状态 待处理/处理中/已解决"`
	HandledBy   string     `json:"handledBy" gorm:"type:varchar(100);comment:处理人"`
	HandledAt   *time.Time `json:"handledAt" gorm:"comment:处理时间"`
	Resolution  string     `json:"resolution" gorm:"type:varchar(300);comment:处理结果"`
}

func (VehicleFault) TableName() string { return "vehicle_faults" }

// BoardingVerification 上车核验（新建表）
type BoardingVerification struct {
	global.GVA_MODEL
	OrderID        uint       `json:"orderId" gorm:"comment:订单ID"`
	PassengerID    uint       `json:"passengerId" gorm:"comment:乘客ID"`
	VehicleID      uint       `json:"vehicleId" gorm:"comment:车辆ID"`
	VerifyType     string     `json:"verifyType" gorm:"type:varchar(50);comment:核验方式 手机尾号/NFC/动态码"`
	VerifyCodeHash string     `json:"-" gorm:"type:varchar(200);comment:凭据哈希(不存明文)"`
	VerifyStatus   string     `json:"verifyStatus" gorm:"type:varchar(50);comment:核验结果 成功/失败"`
	VerifiedAt     *time.Time `json:"verifiedAt" gorm:"comment:核验时间"`
}

func (BoardingVerification) TableName() string { return "boarding_verifications" }

// SafetyIncident 安全事件（字段类型对齐 ride-hailing 实际表结构）
type SafetyIncident struct {
	BizModel
	IncidentNo   string     `json:"incidentNo" gorm:"type:varchar(50);comment:事件编号"`
	OrderID      uint       `json:"orderId" gorm:"comment:关联订单ID"`
	PassengerID  uint       `json:"passengerId" gorm:"comment:乘客ID"`
	IncidentType int        `json:"incidentType" gorm:"type:smallint;comment:事件类型 1紧急制动 2乘客不适 3其他"`
	Description  string     `json:"description" gorm:"type:varchar(500);comment:事件描述"`
	Status       int        `json:"status" gorm:"type:smallint;comment:状态 0待处理 1处理中 2已办结"`
	HandledBy    uint       `json:"handledBy" gorm:"comment:处理人ID"`
	HandledAt    *time.Time `json:"handledAt" gorm:"comment:处理时间"`
	Resolution   string     `json:"resolution" gorm:"type:varchar(500);comment:处理结果"`
}

func (SafetyIncident) TableName() string { return "safety_incidents" }

// ==================== 以下为 PRD 数据字典中存在、库中缺失的表（本次补建） ====================

// OrderCancellation 订单取消记录（模块2）
type OrderCancellation struct {
	global.GVA_MODEL
	OrderID      uint    `json:"orderId" gorm:"comment:订单ID"`
	PassengerID  uint    `json:"passengerId" gorm:"comment:乘客ID"`
	CancelType   string  `json:"cancelType" gorm:"type:varchar(50);comment:取消类型 乘客取消/系统取消"`
	CancelReason string  `json:"cancelReason" gorm:"type:varchar(255);comment:取消原因"`
	CancelFee    float64 `json:"cancelFee" gorm:"type:numeric(10,2);comment:取消费用"`
	Status       string  `json:"status" gorm:"type:varchar(50);comment:状态"`
}

func (OrderCancellation) TableName() string { return "order_cancellations" }

// OrderFeeItem 订单费用明细（模块2）
type OrderFeeItem struct {
	global.GVA_MODEL
	OrderID        uint    `json:"orderId" gorm:"comment:订单ID"`
	FeeType        string  `json:"feeType" gorm:"type:varchar(50);comment:费用类型 起步价/里程费/时长费/附加费/优惠"`
	FeeName        string  `json:"feeName" gorm:"type:varchar(100);comment:费用名称"`
	Amount         float64 `json:"amount" gorm:"type:numeric(10,2);comment:金额"`
	DiscountAmount float64 `json:"discountAmount" gorm:"type:numeric(10,2);comment:优惠金额(负数)"`
}

func (OrderFeeItem) TableName() string { return "order_fee_items" }

// Refund 退款单（模块2）
type Refund struct {
	global.GVA_MODEL
	RefundNo     string     `json:"refundNo" gorm:"type:varchar(50);comment:退款单号"`
	OrderID      uint       `json:"orderId" gorm:"comment:订单ID"`
	PaymentID    uint       `json:"paymentId" gorm:"comment:支付单ID"`
	RefundAmount float64    `json:"refundAmount" gorm:"type:numeric(10,2);comment:退款金额"`
	RefundStatus string     `json:"refundStatus" gorm:"type:varchar(50);comment:状态 待处理/退款中/已退款/失败"`
	RefundReason string     `json:"refundReason" gorm:"type:varchar(255);comment:退款原因"`
	RefundedAt   *time.Time `json:"refundedAt" gorm:"comment:退款时间"`
}

func (Refund) TableName() string { return "refunds" }

// ServiceAreaSchedule 运营区域服务时段（模块3）
type ServiceAreaSchedule struct {
	global.GVA_MODEL
	ServiceAreaID uint       `json:"serviceAreaId" gorm:"comment:运营区域ID"`
	DayOfWeek     int        `json:"dayOfWeek" gorm:"comment:星期 1-7"`
	StartTime     string     `json:"startTime" gorm:"type:varchar(10);comment:开始运营时间 HH:MM"`
	EndTime       string     `json:"endTime" gorm:"type:varchar(10);comment:结束运营时间 HH:MM"`
	IsEnabled     bool       `json:"isEnabled" gorm:"type:boolean;default:true;comment:是否启用"`
	SpecialDate   *time.Time `json:"specialDate" gorm:"comment:特殊日期(为空表示每周规则)"`
	Notice        string     `json:"notice" gorm:"type:varchar(300);comment:时段提示"`
}

func (ServiceAreaSchedule) TableName() string { return "service_area_schedules" }

// VehicleStatusLog 车辆状态日志（模块3）
type VehicleStatusLog struct {
	global.GVA_MODEL
	VehicleID      uint      `json:"vehicleId" gorm:"comment:车辆ID"`
	FromStatus     int       `json:"fromStatus" gorm:"comment:原状态"`
	ToStatus       int       `json:"toStatus" gorm:"comment:新状态"`
	Longitude      float64   `json:"longitude" gorm:"type:numeric(10,6);comment:经度"`
	Latitude       float64   `json:"latitude" gorm:"type:numeric(10,6);comment:纬度"`
	BatteryPercent float64   `json:"batteryPercent" gorm:"type:numeric(5,2);comment:电量百分比"`
	ReportedAt     time.Time `json:"reportedAt" gorm:"comment:上报时间"`
}

func (VehicleStatusLog) TableName() string { return "vehicle_status_logs" }

// VehicleLocationLog 车辆位置日志（模块4）
type VehicleLocationLog struct {
	global.GVA_MODEL
	VehicleID  uint      `json:"vehicleId" gorm:"comment:车辆ID"`
	OrderID    uint      `json:"orderId" gorm:"comment:关联订单ID"`
	Longitude  float64   `json:"longitude" gorm:"type:numeric(10,6);comment:经度"`
	Latitude   float64   `json:"latitude" gorm:"type:numeric(10,6);comment:纬度"`
	SpeedKmh   float64   `json:"speedKmh" gorm:"type:numeric(6,2);comment:速度km/h"`
	RecordedAt time.Time `json:"recordedAt" gorm:"comment:记录时间"`
}

func (VehicleLocationLog) TableName() string { return "vehicle_location_logs" }

// VehicleTelemetryLog 车辆遥测日志（模块4）
type VehicleTelemetryLog struct {
	global.GVA_MODEL
	VehicleID   uint      `json:"vehicleId" gorm:"comment:车辆ID"`
	MetricType  string    `json:"metricType" gorm:"type:varchar(50);comment:指标类型 battery/speed/temperature等"`
	MetricValue float64   `json:"metricValue" gorm:"type:numeric(12,4);comment:指标值"`
	ReportedAt  time.Time `json:"reportedAt" gorm:"comment:上报时间"`
}

func (VehicleTelemetryLog) TableName() string { return "vehicle_telemetry_logs" }

// VehicleMediaRecord 车辆媒体记录（模块4）
type VehicleMediaRecord struct {
	global.GVA_MODEL
	VehicleID  uint      `json:"vehicleId" gorm:"comment:车辆ID"`
	OrderID    uint      `json:"orderId" gorm:"comment:关联订单ID"`
	MediaType  string    `json:"mediaType" gorm:"type:varchar(20);comment:媒体类型 image/video/audio"`
	MediaURL   string    `json:"mediaUrl" gorm:"type:varchar(500);comment:媒体地址"`
	RecordedAt time.Time `json:"recordedAt" gorm:"comment:记录时间"`
}

func (VehicleMediaRecord) TableName() string { return "vehicle_media_records" }
