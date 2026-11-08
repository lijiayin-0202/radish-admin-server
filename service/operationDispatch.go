package service

import (
	"errors"

	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
)

// OperationDispatchService 运营调度服务（bizdb postgres 库 operation_dispatch）
type OperationDispatchService struct{}

func (s *OperationDispatchService) CreateDispatch(a *tripModel.OperationDispatch) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Create(a).Error
}

func (s *OperationDispatchService) UpdateDispatch(a *tripModel.OperationDispatch) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Model(a).Omit("created_at").Updates(a).Error
}

func (s *OperationDispatchService) DeleteDispatch(id uint) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.OperationDispatch{}, id).Error
}

func (s *OperationDispatchService) DeleteDispatchByIds(ids []int) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.OperationDispatch{}, ids).Error
}

func (s *OperationDispatchService) GetDispatch(id uint) (r tripModel.OperationDispatch, err error) {
	d, err := opsBizDB()
	if err != nil {
		return
	}
	err = d.Where("id = ?", id).First(&r).Error
	return
}

// CancelDispatch 撤销派单：将调度状态置为「已取消」(4)。
// 仅允许对未结束（非已完成、非已取消）的调度执行撤销。
func (s *OperationDispatchService) CancelDispatch(id uint) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	var r tripModel.OperationDispatch
	if err = d.Where("id = ?", id).First(&r).Error; err != nil {
		return err
	}
	if r.Status == 3 || r.Status == 4 {
		return errors.New("该调度已结束或已取消，无法撤销")
	}
	return d.Model(&tripModel.OperationDispatch{}).Where("id = ?", id).Update("status", 4).Error
}

func (s *OperationDispatchService) GetDispatchList(info tripReq.OperationDispatchSearch) (list []tripModel.OperationDispatch, total int64, err error) {
	d, err := opsBizDB()
	if err != nil {
		return
	}
	q := d.Model(&tripModel.OperationDispatch{})
	if info.DispatchNo != "" {
		q = q.Where("dispatch_no LIKE ?", "%"+info.DispatchNo+"%")
	}
	if info.DriverID > 0 {
		q = q.Where("driver_id = ?", info.DriverID)
	}
	if info.CarID > 0 {
		q = q.Where("car_id = ?", info.CarID)
	}
	if info.OrderID > 0 {
		q = q.Where("order_id = ?", info.OrderID)
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	if info.StartTime != "" {
		q = q.Where("dispatch_time >= ?", info.StartTime)
	}
	if info.EndTime != "" {
		q = q.Where("dispatch_time <= ?", info.EndTime)
	}
	if err = q.Count(&total).Error; err != nil {
		return
	}
	page, pageSize := info.Page, info.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	err = q.Limit(pageSize).Offset(pageSize * (page - 1)).Order("id desc").Find(&list).Error
	return

}
