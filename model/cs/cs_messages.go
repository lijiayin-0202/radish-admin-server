package cs

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CsMessage 客服会话消息表，同一会话内按 seq 递增，按 seq 正序即为会话历史消息
// 说明：软删除统一使用 global.GVA_MODEL 自带的 DeletedAt，本结构体不再声明删除相关字段
type CsMessage struct {
	global.GVA_MODEL
	SessionID  int64      `json:"sessionId" form:"sessionId" gorm:"column:session_id;type:bigint;not null;default:0;comment:归属会话ID，关联cs_sessions"`
	SenderType int64      `json:"senderType" form:"senderType" gorm:"column:sender_type;not null;default:1;comment:发送方类型：1乘客，2客服，3系统"`
	MsgType    int64      `json:"msgType" form:"msgType" gorm:"column:msg_type;not null;default:1;comment:消息类型：1文本，2图片，3文件，4系统通知"`
	Content    string     `json:"content" form:"content" gorm:"column:content;type:text;not null;default:'';comment:消息文本内容"`
	AttachURL  string     `json:"attachUrl" form:"attachUrl" gorm:"column:attach_url;type:varchar(500);not null;default:'';comment:附件地址"`
	Seq        int64      `json:"seq" form:"seq" gorm:"column:seq;type:bigint;not null;default:0;comment:消息序号，同一会话内递增"`
	IsRead     int64      `json:"isRead" form:"isRead" gorm:"column:is_read;not null;default:0;comment:已读状态：0未读，1已读"`
	SendTime   *time.Time `json:"sendTime" form:"sendTime" gorm:"column:send_time;comment:消息发送时间"`
}

// TableName CsMessage 自定义表名 cs_messages
func (CsMessage) TableName() string {
	return "cs_messages"
}
