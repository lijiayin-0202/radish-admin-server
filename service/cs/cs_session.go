package cs

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
)

type CsSessionService struct{}

// 客服会话状态：1进行中，2已结束（取值与 cs_sessions.status 注释一致）
const (
	CsSessionStatusOngoing int64 = 1 // 进行中
	CsSessionStatusEnded   int64 = 2 // 已结束
)

// csClerkAuthorityName 客服角色名，对应系统库 sys_authorities.authority_name
// 项目未内置客服角色，需后台「角色管理」创建同名角色并把客服账号授权到该角色
const csClerkAuthorityName = "客服"

// csClerkAssignCursor 轮询游标：记录下一次分配的客服下标，配合互斥锁保证并发下同一客服轮流接待
// 说明：进程内内存计数，服务重启后归零，适用于单实例部署
var (
	csClerkAssignMu     sync.Mutex
	csClerkAssignCursor int
)

// checkCsSessionStatus 校验会话状态取值
func checkCsSessionStatus(status int64) error {
	if status != CsSessionStatusOngoing && status != CsSessionStatusEnded {
		return errors.New("会话状态取值不合法，仅支持：1进行中，2已结束")
	}
	return nil
}

// GetCsSessionList 分页获取客服会话列表，支持状态、客服、用户与关键词筛选
func (s *CsSessionService) GetCsSessionList(search csReq.CsSessionSearch) (list []csModel.CsSession, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsSession{})
	if search.SessionNo != "" {
		db = db.Where("session_no = ?", search.SessionNo)
	}
	if search.Keyword != "" {
		kw := "%" + search.Keyword + "%"
		db = db.Where("session_no LIKE ? OR user_name LIKE ? OR clerk_name LIKE ?", kw, kw, kw)
	}
	if search.TicketID != nil {
		db = db.Where("ticket_id = ?", *search.TicketID)
	}
	if search.UserID != nil {
		db = db.Where("user_id = ?", *search.UserID)
	}
	if search.UserPhone != "" {
		db = db.Where("user_phone = ?", search.UserPhone)
	}
	if search.ClerkID != nil {
		db = db.Where("clerk_id = ?", *search.ClerkID)
	}
	if search.Status != nil {
		db = db.Where("status = ?", *search.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// listOnlineCsClerks 查询在线客服池
// 数据源为系统库 sys_users：取未冻结（enable=1）且属于客服角色的用户，按用户ID升序保证轮询顺序稳定
// 说明：GVA 不维护用户在线状态，此处以「未冻结 + 客服角色」近似在线
func listOnlineCsClerks() (clerks []sysModel.SysUser, err error) {
	if global.GVA_DB == nil {
		return nil, errors.New("系统库未初始化，无法查询客服列表")
	}
	var authorityIDs []uint
	if err = global.GVA_DB.Model(&sysModel.SysAuthority{}).
		Where("authority_name = ?", csClerkAuthorityName).
		Pluck("authority_id", &authorityIDs).Error; err != nil {
		global.GVA_LOG.Error("查询客服角色失败!", zap.String("authorityName", csClerkAuthorityName), zap.Error(err))
		return nil, errors.New("查询客服角色失败")
	}
	if len(authorityIDs) == 0 {
		return nil, errors.New("后台未创建客服角色，无法分配客服")
	}
	if err = global.GVA_DB.Model(&sysModel.SysUser{}).
		Select("id", "username", "nick_name").
		Where("enable = ?", 1).
		Where("authority_id IN ?", authorityIDs).
		Order("id ASC").
		Find(&clerks).Error; err != nil {
		global.GVA_LOG.Error("查询在线客服列表失败!", zap.Error(err))
		return nil, errors.New("查询在线客服列表失败")
	}
	return clerks, nil
}

// pickCsClerkByRoundRobin 在客服池中轮询选取一位客服
// 进程内游标 + 互斥锁：并发创建会话时同一批客服轮流接待，不会连续分给同一人
func pickCsClerkByRoundRobin(clerks []sysModel.SysUser) sysModel.SysUser {
	csClerkAssignMu.Lock()
	defer csClerkAssignMu.Unlock()
	idx := csClerkAssignCursor % len(clerks)
	csClerkAssignCursor = (idx + 1) % len(clerks)
	return clerks[idx]
}

// assignCsClerk 自动分配接待客服，返回客服ID与客服姓名
// 客服池为空或查询失败时返回业务错误，由调用方决定是否继续创建会话
func assignCsClerk() (clerkID int64, clerkName string, err error) {
	clerks, err := listOnlineCsClerks()
	if err != nil {
		return 0, "", err
	}
	if len(clerks) == 0 {
		return 0, "", errors.New("暂无在线客服，无法创建会话")
	}
	clerk := pickCsClerkByRoundRobin(clerks)
	clerkName = clerk.NickName
	if clerkName == "" {
		clerkName = clerk.Username // 昵称未填时退回登录名，避免落库后客服姓名为空
	}
	global.GVA_LOG.Info("会话自动分配客服", zap.Uint("clerkId", clerk.ID), zap.String("clerkName", clerkName))
	return int64(clerk.ID), clerkName, nil
}

// CreateCsSession 新增客服会话
// 会话编号由 global.Snow 雪花ID自动生成，前端不传；新建会话固定为进行中，开始时间为空时取服务器当前时间
// 入参未指定 clerkId 时按轮询自动分配在线客服，写入 clerk_id / clerk_name；无可用客服则会话不创建
// 返回落库后的会话实体，便于前端拿到 sessionId / sessionNo 继续发送会话消息
func (s *CsSessionService) CreateCsSession(operatorID uint, req csReq.CsSessionCreateReq) (entity csModel.CsSession, err error) {
	startedAt := req.StartedAt
	if startedAt == nil {
		now := time.Now()
		startedAt = &now
	}
	entity = csModel.CsSession{
		SessionNo: strconv.FormatInt(global.Snow.NextID(), 10), // 会话编号后端自动生成
		UserID:    req.UserID,
		UserName:  req.UserName,
		UserPhone: req.UserPhone,
		ClerkName: req.ClerkName,
		Status:    CsSessionStatusOngoing, // 新建会话固定为进行中
		StartedAt: startedAt,
		Remark:    req.Remark,
	}
	if req.TicketID != nil {
		entity.TicketID = *req.TicketID
	}
	if req.ClerkID != nil {
		entity.ClerkID = *req.ClerkID // 入参指定客服时按指定分配，不消耗轮询游标
	} else {
		// 未指定客服：轮询自动分配；无可用客服直接返回业务错误，会话不落库
		clerkID, clerkName, assignErr := assignCsClerk()
		if assignErr != nil {
			return entity, assignErr
		}
		entity.ClerkID = clerkID
		entity.ClerkName = clerkName
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsSession 修改客服会话，会话编号不允许修改
// Status 传 1 进行中 / 2 已结束，nil 表示不变更状态；流转到已结束时结束时间为空则由后端补当前时间
func (s *CsSessionService) UpdateCsSession(operatorID uint, req csReq.CsSessionUpdateReq) (err error) {
	var entity csModel.CsSession
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("会话不存在或已删除")
	}
	if req.Status != nil {
		if err = checkCsSessionStatus(*req.Status); err != nil {
			return
		}
		if *req.Status == CsSessionStatusEnded && req.EndedAt == nil && entity.EndedAt == nil {
			now := time.Now()
			req.EndedAt = &now
		}
		entity.Status = *req.Status
	}

	entity.UserID = req.UserID
	entity.UserName = req.UserName
	entity.UserPhone = req.UserPhone
	entity.ClerkName = req.ClerkName
	entity.Remark = req.Remark
	if req.TicketID != nil {
		entity.TicketID = *req.TicketID
	}
	if req.ClerkID != nil {
		entity.ClerkID = *req.ClerkID
	}
	if req.StartedAt != nil {
		entity.StartedAt = req.StartedAt // 开始时间空值不覆盖已有内容
	}
	if req.EndedAt != nil {
		entity.EndedAt = req.EndedAt // 结束时间空值不覆盖已有内容
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsSession 删除客服会话（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsSessionService) DeleteCsSession(id uint) (err error) {
	var entity csModel.CsSession
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("会话不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsSession 获取客服会话详情
func (s *CsSessionService) GetCsSession(id uint) (entity csModel.CsSession, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
