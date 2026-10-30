package cs

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsFaqService struct{}

// FAQ 上下架状态：0禁用，1启用（取值与 cs_faqs.status 注释一致）
const (
	CsFaqStatusOff int64 = 0 // 下架（禁用）
	CsFaqStatusOn  int64 = 1 // 上架（启用）
)

// checkCsFaqStatus 校验 FAQ 上架状态取值
func checkCsFaqStatus(status int64) error {
	if status != CsFaqStatusOff && status != CsFaqStatusOn {
		return errors.New("FAQ状态取值不合法，仅支持：0禁用，1启用")
	}
	return nil
}

// GetCsFaqList 分页获取常见问题列表，支持分类、上架状态与关键词筛选
func (s *CsFaqService) GetCsFaqList(search csReq.CsFaqSearch) (list []csModel.CsFaq, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsFaq{})
	if search.Question != "" {
		db = db.Where("question LIKE ?", "%"+search.Question+"%")
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("question LIKE ? OR answer LIKE ? OR tags LIKE ?", kw, kw, kw)
	}
	if search.Category != "" {
		db = db.Where("category LIKE ?", "%"+search.Category+"%")
	}
	if search.Tags != "" {
		db = db.Where("tags LIKE ?", "%"+search.Tags+"%")
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("sort DESC").Order("id DESC").Find(&list).Error
	return
}

// csFaqKeywordEscape LIKE 通配符转义器，防止用户输入的 %、_ 被数据库当成通配符
// 反斜杠自身必须先转义：否则关键词末尾的 \ 会吃掉后面的转义符，语义被反转
var csFaqKeywordEscape = strings.NewReplacer(
	`\`, `\\`,
	`%`, `\%`,
	`_`, `\_`,
)

// SearchCsFaqByKeyword 按关键词模糊检索常见问题
// 对问题标题（question）与回答内容（answer）做 LIKE 模糊匹配，复用现有分页结构与列表排序规则
// 关键词中的 LIKE 通配符（% / _ / \）先转义，再用 ESCAPE '\' 声明转义符，适配 PostgreSQL
// 关键词为空或纯空白时直接返回空数组，避免无条件下发全表查询
func (s *CsFaqService) SearchCsFaqByKeyword(search csReq.CsFaqSearch) (list []csModel.CsFaq, total int64, err error) {
	list = make([]csModel.CsFaq, 0) // 保证空结果返回 [] 而不是 null
	keyword := strings.TrimSpace(search.Keyword)
	if keyword == "" {
		return list, 0, nil
	}
	kw := "%" + csFaqKeywordEscape.Replace(keyword) + "%"
	db := global.LUOBO_DB.Model(&csModel.CsFaq{}).
		Where("question LIKE ? OR answer LIKE ? ESCAPE '\\'", kw, kw)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("sort DESC").Order("id DESC").Find(&list).Error
	return
}

// CreateCsFaq 新增常见问题
// Status 传 0 禁用 / 1 启用，nil 表示按库默认1启用
func (s *CsFaqService) CreateCsFaq(operatorID uint, req csReq.CsFaqCreateReq) (entity csModel.CsFaq, err error) {
	status := req.Status
	if status == nil {
		on := CsFaqStatusOn
		status = &on
	}
	if err = checkCsFaqStatus(*status); err != nil {
		return entity, err
	}
	entity = csModel.CsFaq{
		Question: req.Question,
		Answer:   req.Answer,
		Category: req.Category,
		Tags:     req.Tags,
		Sort:     req.Sort,
		Status:   *status,
		Remark:   req.Remark,
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsFaq 修改常见问题，含上下架：Status 传 0 下架 / 1 上架，nil 表示不变更上架状态
func (s *CsFaqService) UpdateCsFaq(operatorID uint, req csReq.CsFaqUpdateReq) (err error) {
	var entity csModel.CsFaq
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("FAQ不存在或已删除")
	}
	entity.Question = req.Question
	entity.Answer = req.Answer
	entity.Category = req.Category
	entity.Tags = req.Tags
	entity.Sort = req.Sort
	entity.Remark = req.Remark
	if req.Status != nil {
		if err = checkCsFaqStatus(*req.Status); err != nil {
			return
		}
		entity.Status = *req.Status
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsFaq 删除常见问题（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsFaqService) DeleteCsFaq(id uint) (err error) {
	var entity csModel.CsFaq
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("FAQ不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsFaq 获取常见问题详情
func (s *CsFaqService) GetCsFaq(id uint) (entity csModel.CsFaq, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
