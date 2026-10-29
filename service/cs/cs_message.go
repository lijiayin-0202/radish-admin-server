package cs

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
	csReq "github.com/flipped-aurora/gin-vue-admin/server/model/cs/request"
)

type CsMessageService struct{}

// 消息发送方类型 / 消息类型，取值与 cs_messages 注释一致
const (
	CsMessageSenderPassenger int64 = 1 // 发送方：乘客
	CsMessageSenderClerk     int64 = 2 // 发送方：客服
	CsMessageSenderSystem    int64 = 3 // 发送方：系统

	CsMessageTypeText    int64 = 1 // 消息类型：文本
	CsMessageTypeImage   int64 = 2 // 消息类型：图片
	CsMessageTypeFile    int64 = 3 // 消息类型：文件
	CsMessageTypeSystem  int64 = 4 // 消息类型：系统通知

	CsMessageUnread int64 = 0 // 已读状态：未读
	CsMessageRead   int64 = 1 // 已读状态：已读
)

// checkCsMessageSenderType 校验消息发送方类型
func checkCsMessageSenderType(senderType int64) error {
	if senderType < CsMessageSenderPassenger || senderType > CsMessageSenderSystem {
		return errors.New("发送方类型取值不合法，仅支持：1乘客，2客服，3系统")
	}
	return nil
}

// checkCsMessageMsgType 校验消息类型
func checkCsMessageMsgType(msgType int64) error {
	if msgType < CsMessageTypeText || msgType > CsMessageTypeSystem {
		return errors.New("消息类型取值不合法，仅支持：1文本，2图片，3文件，4系统通知")
	}
	return nil
}

// GetCsMessageList 分页获取客服会话消息列表
// 传 SessionID 时按会话历史消息语义返回（消息序号正序，再按主键正序），否则按最新消息倒序返回
func (s *CsMessageService) GetCsMessageList(search csReq.CsMessageSearch) (list []csModel.CsMessage, total int64, err error) {
	db := global.LUOBO_DB.Model(&csModel.CsMessage{})
	if search.SessionID != nil {
		db = db.Where("session_id = ?", *search.SessionID)
	}
	if search.SenderType != nil {
		db = db.Where("sender_type = ?", *search.SenderType)
	}
	if search.MsgType != nil {
		db = db.Where("msg_type = ?", *search.MsgType)
	}
	if search.IsRead != nil {
		db = db.Where("is_read = ?", *search.IsRead)
	}
	if search.Keyword != "" {
		db = db.Where("content LIKE ?", "%"+search.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if search.SessionID != nil {
		err = db.Scopes(search.PageInfo.Paginate()).Order("seq ASC").Order("id ASC").Find(&list).Error
		return
	}
	err = db.Scopes(search.PageInfo.Paginate()).Order("id DESC").Find(&list).Error
	return
}

// CreateCsMessage 保存客服会话消息
// 校验归属会话存在；消息序号不传时按当前会话最大序号+1自动生成；消息发送时间不传时取服务器当前时间
// 返回落库后的消息实体
// TODO(ride-hailing)：消息如需附带订单/车辆快照，调用 ride-hailing 微服务「订单查询接口」获取，
// 本仓库不落地订单/车辆 model
func (s *CsMessageService) CreateCsMessage(operatorID uint, req csReq.CsMessageCreateReq) (entity csModel.CsMessage, err error) {
	var session csModel.CsSession
	if err = global.LUOBO_DB.Where("id = ?", req.SessionID).First(&session).Error; err != nil {
		return entity, errors.New("归属会话不存在或已删除")
	}
	if err = checkCsMessageSenderType(req.SenderType); err != nil {
		return entity, err
	}
	if err = checkCsMessageMsgType(req.MsgType); err != nil {
		return entity, err
	}

	seq := req.Seq
	if seq <= 0 {
		var maxSeq int64
		if err = global.LUOBO_DB.Model(&csModel.CsMessage{}).Where("session_id = ?", req.SessionID).
			Select("COALESCE(MAX(seq), 0)").Scan(&maxSeq).Error; err != nil {
			return entity, errors.New("生成消息序号失败")
		}
		seq = maxSeq + 1
	}
	sendTime := req.SendTime
	if sendTime == nil {
		now := time.Now()
		sendTime = &now
	}

	entity = csModel.CsMessage{
		SessionID:  req.SessionID,
		SenderType: req.SenderType,
		MsgType:    req.MsgType,
		Content:    req.Content,
		AttachURL:  req.AttachUrl,
		Seq:        seq,
		IsRead:     req.IsRead,
		SendTime:   sendTime,
	}
	err = global.LUOBO_DB.Create(&entity).Error
	return
}

// UpdateCsMessage 修改客服会话消息
func (s *CsMessageService) UpdateCsMessage(operatorID uint, req csReq.CsMessageUpdateReq) (err error) {
	var entity csModel.CsMessage
	if err = global.LUOBO_DB.Where("id = ?", req.ID).First(&entity).Error; err != nil {
		return errors.New("消息不存在或已删除")
	}
	if err = checkCsMessageSenderType(req.SenderType); err != nil {
		return
	}
	if err = checkCsMessageMsgType(req.MsgType); err != nil {
		return
	}
	entity.SenderType = req.SenderType
	entity.MsgType = req.MsgType
	entity.Content = req.Content
	entity.AttachURL = req.AttachUrl
	entity.Seq = req.Seq
	entity.IsRead = req.IsRead
	if req.SendTime != nil {
		entity.SendTime = req.SendTime // 消息发送时间空值不覆盖已有内容
	}
	return global.LUOBO_DB.Save(&entity).Error
}

// DeleteCsMessage 删除客服会话消息（软删除，依赖 global.GVA_MODEL 自带的 deleted_at）
func (s *CsMessageService) DeleteCsMessage(id uint) (err error) {
	var entity csModel.CsMessage
	if err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error; err != nil {
		return errors.New("消息不存在或已删除")
	}
	return global.LUOBO_DB.Delete(&entity).Error
}

// GetCsMessage 获取客服会话消息详情
func (s *CsMessageService) GetCsMessage(id uint) (entity csModel.CsMessage, err error) {
	err = global.LUOBO_DB.Where("id = ?", id).First(&entity).Error
	return
}
