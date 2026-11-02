package trip

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/service/bizdb"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	"gorm.io/gorm"
)

type TripService struct{}

func db() (*gorm.DB, error) {
	d := bizdb.PostgresDB()
	if d == nil {
		return nil, errors.New("出行业务库连接不可用，请检查数据库配置")
	}
	return d, nil
}

// paginate 通用分页
func paginate[T any](d *gorm.DB, page, pageSize int, out *[]T) (int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var total int64
	if err := d.Count(&total).Error; err != nil {
		return 0, err
	}
	err := d.Limit(pageSize).Offset(pageSize * (page - 1)).Order("id desc").Find(out).Error
	return total, err
}

// scopeAreaByServiceAreaID 按 service_area_id 过滤城市/区域（order/station/vehicle 等直接带该字段的表）
func scopeAreaByServiceAreaID(q *gorm.DB, area uint, city string) *gorm.DB {
	if area > 0 {
		return q.Where("service_area_id = ?", area)
	}
	if city != "" {
		return q.Where("service_area_id IN (SELECT id FROM service_areas WHERE city_code = ?)", city)
	}
	return q
}

// scopeAreaByVehicle 通过 vehicle_id 关联 vehicles 的 service_area_id 过滤城市/区域（派单/指令/故障/核验等按车辆维度归属区域）
func scopeAreaByVehicle(q *gorm.DB, area uint, city string) *gorm.DB {
	if area > 0 {
		return q.Where("vehicle_id IN (SELECT id FROM vehicles WHERE service_area_id = ?)", area)
	}
	if city != "" {
		return q.Where("vehicle_id IN (SELECT id FROM vehicles WHERE service_area_id IN (SELECT id FROM service_areas WHERE city_code = ?))", city)
	}
	return q
}

// scopeTime 按时间列过滤范围（col 为该实体的时间字段，如 created_at / requested_at）
func scopeTime(q *gorm.DB, startTime, endTime, col string) *gorm.DB {
	if startTime != "" {
		q = q.Where(col+" >= ?", startTime)
	}
	if endTime != "" {
		q = q.Where(col+" <= ?", endTime)
	}
	return q
}

// ==================== 乘客 ====================

func (s *TripService) GetPassengerList(info tripReq.TripSearch) (list []tripModel.Passenger, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.Passenger{})
	if info.Keyword != "" {
		q = q.Where("mobile LIKE ? OR nickname LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at") // 乘客无区域维度，仅按时间过滤
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreatePassenger(p *tripModel.Passenger) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(p).Error
}

func (s *TripService) UpdatePassenger(p *tripModel.Passenger) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(p).Omit("created_at").Updates(p).Error
}

func (s *TripService) DeletePassenger(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.Passenger{}, id).Error
}

// ==================== 订单 ====================

func (s *TripService) GetOrderList(info tripReq.OrderSearch) (list []tripModel.TripOrder, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.TripOrder{})
	if info.Keyword != "" {
		q = q.Where("order_no LIKE ? OR pickup_name LIKE ? OR destination_name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.OrderStatus != nil {
		q = q.Where("order_status = ?", *info.OrderStatus)
	}
	q = scopeAreaByServiceAreaID(q, info.Area, info.City) // 订单按 service_area_id 过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "requested_at") // 订单时间用下单时间
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) GetOrderStatusLogs(orderID uint) (list []tripModel.TripStatusLog, err error) {
	d, err := db()
	if err != nil {
		return
	}
	err = d.Where("order_id = ?", orderID).Order("id asc").Find(&list).Error
	return
}

// ==================== 运营区域 ====================

func (s *TripService) GetAreaList(info tripReq.TripSearch) (list []tripModel.ServiceArea, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.ServiceArea{})
	if info.Keyword != "" {
		q = q.Where("area_name LIKE ? OR city_name LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	if info.City != "" {
		q = q.Where("city_code = ?", info.City) // 区域页按城市编码过滤
	}
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreateArea(a *tripModel.ServiceArea) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(a).Error
}

func (s *TripService) UpdateArea(a *tripModel.ServiceArea) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(a).Omit("created_at").Updates(a).Error
}

func (s *TripService) DeleteArea(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.ServiceArea{}, id).Error
}

// ==================== 站点 ====================

func (s *TripService) GetStationList(info tripReq.TripSearch) (list []tripModel.Station, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.Station{})
	if info.Keyword != "" {
		q = q.Where("station_name LIKE ? OR address LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	q = scopeAreaByServiceAreaID(q, info.Area, info.City) // 站点按 service_area_id 过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreateStation(a *tripModel.Station) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(a).Error
}

func (s *TripService) UpdateStation(a *tripModel.Station) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(a).Omit("created_at").Updates(a).Error
}

func (s *TripService) DeleteStation(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.Station{}, id).Error
}

// ==================== 车辆 ====================

func (s *TripService) GetVehicleList(info tripReq.TripSearch) (list []tripModel.Vehicle, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.Vehicle{})
	if info.Keyword != "" {
		q = q.Where("vehicle_no LIKE ? OR plate_no LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	q = scopeAreaByServiceAreaID(q, info.Area, info.City) // 车辆按 service_area_id 过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreateVehicle(v *tripModel.Vehicle) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(v).Error
}

func (s *TripService) UpdateVehicle(v *tripModel.Vehicle) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(v).Omit("created_at").Updates(v).Error
}

func (s *TripService) DeleteVehicle(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.Vehicle{}, id).Error
}

// ==================== 派单记录 ====================

func (s *TripService) GetDispatchList(info tripReq.TripSearch) (list []tripModel.DispatchRecord, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.DispatchRecord{})
	if info.Keyword != "" {
		q = q.Where("dispatch_no LIKE ? OR dispatch_status LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	q = scopeAreaByVehicle(q, info.Area, info.City) // 派单按车辆归属区域过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

// ==================== 车辆指令 ====================

func (s *TripService) GetCommandList(info tripReq.TripSearch) (list []tripModel.VehicleCommand, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.VehicleCommand{})
	if info.Keyword != "" {
		q = q.Where("command_no LIKE ? OR command_type LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	q = scopeAreaByVehicle(q, info.Area, info.City) // 指令按车辆归属区域过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreateCommand(cmd *tripModel.VehicleCommand) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(cmd).Error
}

func (s *TripService) UpdateCommand(cmd *tripModel.VehicleCommand) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(cmd).Omit("created_at").Updates(cmd).Error
}

// ==================== 车辆故障 ====================

func (s *TripService) GetFaultList(info tripReq.TripSearch) (list []tripModel.VehicleFault, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.VehicleFault{})
	if info.Keyword != "" {
		q = q.Where("fault_no LIKE ? OR fault_type LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	q = scopeAreaByVehicle(q, info.Area, info.City) // 故障按车辆归属区域过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) CreateFault(f *tripModel.VehicleFault) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(f).Error
}

func (s *TripService) UpdateFault(f *tripModel.VehicleFault) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(f).Omit("created_at").Updates(f).Error
}

// ==================== 上车核验 ====================

func (s *TripService) GetBoardingList(info tripReq.TripSearch) (list []tripModel.BoardingVerification, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.BoardingVerification{})
	if info.Keyword != "" {
		q = q.Where("verify_type LIKE ? OR verify_status LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	q = scopeAreaByVehicle(q, info.Area, info.City) // 核验按车辆归属区域过滤城市/区域
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at")
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

// ==================== 安全事件 ====================

func (s *TripService) GetIncidentList(info tripReq.TripSearch) (list []tripModel.SafetyIncident, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.SafetyIncident{})
	if info.Keyword != "" {
		q = q.Where("incident_no LIKE ? OR description LIKE ? OR resolution LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	q = scopeTime(q, info.StartTime, info.EndTime, "created_at") // 安全事件无区域维度，仅按时间过滤
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *TripService) UpdateIncident(i *tripModel.SafetyIncident) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(i).Omit("created_at").Updates(i).Error
}
