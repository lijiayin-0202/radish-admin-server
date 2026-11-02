package cs

import (
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/service/bizdb"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	"gorm.io/gorm"
)

type CsService struct{}

func db() (*gorm.DB, error) {
	d := bizdb.LuoboDB()
	if d == nil {
		return nil, errors.New("客服业务库连接不可用，请检查数据库配置")
	}
	return d, nil
}

func paginate[T any](d *gorm.DB, page, pageSize int, out *[]T) (int64, error) {
	var total int64
	if err := d.Count(&total).Error; err != nil {
		return 0, err
	}
	err := d.Limit(pageSize).Offset(pageSize * (page - 1)).Order("id desc").Find(out).Error
	return total, err
}

// ==================== 会话 ====================

func (s *CsService) GetSessionList(info csReq.CsSearch) (list []csModel.CsSession, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&csModel.CsSession{})
	if info.Keyword != "" {
		q = q.Where("session_no LIKE ? OR issue_category LIKE ? OR assigned_to LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

// ==================== 工单 ====================

func (s *CsService) GetTicketList(info csReq.TicketSearch) (list []csModel.CsTicket, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&csModel.CsTicket{})
	if info.Keyword != "" {
		q = q.Where("ticket_no LIKE ? OR title LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *CsService) UpdateTicket(t *csModel.CsTicket) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(t).Omit("created_at").Updates(t).Error
}

// ==================== 工单（新增接口） ====================

func (s *CsService) CreateTicket(t *csModel.CsTicket) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(t).Error
}

func (s *CsService) GetTicket(id uint) (*csModel.CsTicket, error) {
	d, err := db()
	if err != nil {
		return nil, err
	}
	var t csModel.CsTicket
	if err := d.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *CsService) UpdateTicketStatus(id, status int) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(&csModel.CsTicket{}).Where("id = ?", id).Update("status", status).Error
}

// ==================== FAQ 分类 ====================

func (s *CsService) GetFaqCategoryList(info csReq.CsSearch) (list []csModel.CsFaqCategory, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&csModel.CsFaqCategory{})
	if info.Keyword != "" {
		q = q.Where("category_name LIKE ?", "%"+info.Keyword+"%")
	}
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *CsService) CreateFaqCategory(c *csModel.CsFaqCategory) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(c).Error
}

func (s *CsService) UpdateFaqCategory(c *csModel.CsFaqCategory) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(c).Omit("created_at").Updates(c).Error
}

func (s *CsService) DeleteFaqCategory(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&csModel.CsFaqCategory{}, id).Error
}

// ==================== FAQ ====================

func (s *CsService) GetFaqList(info csReq.CsSearch) (list []csModel.CsFaq, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&csModel.CsFaq{})
	if info.Keyword != "" {
		q = q.Where("question LIKE ? OR keywords LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != nil {
		q = q.Where("status = ?", *info.Status)
	}
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}

func (s *CsService) CreateFaq(f *csModel.CsFaq) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Create(f).Error
}

func (s *CsService) UpdateFaq(f *csModel.CsFaq) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Model(f).Omit("created_at").Updates(f).Error
}

func (s *CsService) DeleteFaq(id uint) error {
	d, err := db()
	if err != nil {
		return err
	}
	return d.Delete(&csModel.CsFaq{}, id).Error
}

// ==================== 消息通知 ====================

func (s *CsService) GetNotificationList(info csReq.CsSearch) (list []csModel.CsNotification, total int64, err error) {
	d, err := db()
	if err != nil {
		return
	}
	q := d.Model(&csModel.CsNotification{})
	if info.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+info.Keyword+"%")
	}
	total, err = paginate(q, info.Page, info.PageSize, &list)
	return
}
