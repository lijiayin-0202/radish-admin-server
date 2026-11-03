package service

import (
	"errors"

	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/bizdb"
	"gorm.io/gorm"
)

// ServiceAreaService 运营区域服务（bizdb postgres 库 service_areas）
type ServiceAreaService struct{}

// opsBizDB 获取出行业务库连接
func opsBizDB() (*gorm.DB, error) {
	d := bizdb.PostgresDB()
	if d == nil {
		return nil, errors.New("出行业务库连接不可用，请检查数据库配置")
	}
	return d, nil
}

func (s *ServiceAreaService) CreateServiceArea(a *tripModel.ServiceArea) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Create(a).Error
}

func (s *ServiceAreaService) UpdateServiceArea(a *tripModel.ServiceArea) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Model(a).Omit("created_at").Updates(a).Error
}

func (s *ServiceAreaService) DeleteServiceArea(id uint) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.ServiceArea{}, id).Error
}

func (s *ServiceAreaService) DeleteServiceAreaByIds(ids []int) error {
	d, err := opsBizDB()
	if err != nil {
		return err
	}
	return d.Delete(&tripModel.ServiceArea{}, ids).Error
}

func (s *ServiceAreaService) GetServiceArea(id uint) (r tripModel.ServiceArea, err error) {
	d, err := opsBizDB()
	if err != nil {
		return
	}
	err = d.Where("id = ?", id).First(&r).Error
	return
}

func (s *ServiceAreaService) GetServiceAreaList(info tripReq.TripSearch) (list []tripModel.ServiceArea, total int64, err error) {
	d, err := opsBizDB()
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
		q = q.Where("city_code = ?", info.City)
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
