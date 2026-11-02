package service

import (
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
)

// StationService 站点管理服务（bizdb postgres 库 stations）
type StationService struct{}

func (s *StationService) CreateStation(a *tripModel.Station) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Create(a).Error
}

func (s *StationService) UpdateStation(a *tripModel.Station) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Model(a).Omit("created_at").Updates(a).Error
}

func (s *StationService) DeleteStation(id uint) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.Station{}, id).Error
}

func (s *StationService) DeleteStationByIds(ids []int) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.Station{}, ids).Error
}

func (s *StationService) GetStation(id uint) (r tripModel.Station, err error) {
	d, err := opsBizDB()
	if err != nil {
		return
	}
	err = d.Where("id = ?", id).First(&r).Error
	return
}

func (s *StationService) GetStationList(info tripReq.TripSearch) (list []tripModel.Station, total int64, err error) {
	d, err := opsBizDB()
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
	if info.Area > 0 {
		q = q.Where("service_area_id = ?", info.Area)
	}
	if info.City != "" {
		q = q.Where("service_area_id IN (SELECT id FROM service_areas WHERE city_code = ?)", info.City)
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
