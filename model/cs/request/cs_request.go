package request

import (
	"time"

	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ==========================================================
// cs_tickets 客服工单
// ==========================================================

// CsTicketSearch 客服工单分页查询入参
type CsTicketSearch struct {
	commonReq.PageInfo
	TicketNo          string  `json:"ticketNo" form:"ticketNo"`                   // 工单号，模糊匹配
	Keyword           string  `json:"keyword" form:"keyword"`                     // 综合关键词，工单号/工单标题/问题描述 模糊匹配
	UserID            *int64  `json:"userId" form:"userId"`                       // 用户ID
	UserName          string  `json:"userName" form:"userName"`                   // 用户姓名，模糊匹配
	UserPhone         string  `json:"userPhone" form:"userPhone"`                 // 用户手机号，精确匹配
	OrderID           *int64  `json:"orderId" form:"orderId"`                     // 关联订单ID
	OrderNo           string  `json:"orderNo" form:"orderNo"`                     // 关联订单号，精确匹配
	Status            *int64  `json:"status" form:"status"`                       // 工单状态：1待处理，2处理中，3已解决，4已关闭
	Priority          *int64  `json:"priority" form:"priority"`                 // 优先级：1低，2中，3高
	Category          string  `json:"category" form:"category"`                 // 问题分类，模糊匹配
	AssignedClerkID   *int64  `json:"assignedClerkId" form:"assignedClerkId"`     // 分配客服ID
}

// CsTicketCreateReq 新增客服工单入参
// 工单号由后端 global.Snow 雪花ID自动生成，前端不传；新建工单固定为待受理
// TODO(ride-hailing)：order_id / order_no 只做关联存储，订单与车辆信息调用 ride-hailing 微服务「订单查询接口」获取
type CsTicketCreateReq struct {
	UserID            int64  `json:"userId"`            // 用户ID
	UserName          string `json:"userName"`          // 用户姓名
	UserPhone         string `json:"userPhone"`         // 用户手机号
	OrderID           *int64 `json:"orderId"`           // 关联订单ID，无订单可不传；传了则触发 ride-hailing 订单查询并回填订单号/乘客手机号
	OrderNo           string `json:"orderNo"`           // 关联订单号，无订单可不传
	Priority          int64  `json:"priority"`          // 优先级：1低，2中，3高，不传按库默认1
	Category          string `json:"category"`          // 问题分类
	Title             string `json:"title"`             // 工单标题
	Description       string `json:"description"`       // 问题描述
	AttachUrls        string `json:"attachUrls"`        // 附件URL列表
	AssignedClerkID   *int64 `json:"assignedClerkId"`   // 分配客服ID
	AssignedClerkName string `json:"assignedClerkName"` // 分配客服姓名
	Remark            string `json:"remark"`            // 备注
}

// CsTicketUpdateReq 修改客服工单入参，工单号不允许修改
// Status 使用指针：传值表示按流转规则变更状态（待受理 -> 处理中 -> 已完结），nil 表示不变更状态
type CsTicketUpdateReq struct {
	ID                uint       `json:"id"`                // 主键ID
	UserID            int64      `json:"userId"`            // 用户ID
	UserName          string     `json:"userName"`          // 用户姓名
	UserPhone         string     `json:"userPhone"`         // 用户手机号
	OrderID           *int64     `json:"orderId"`           // 关联订单ID
	OrderNo           string     `json:"orderNo"`           // 关联订单号
	Priority          int64      `json:"priority"`          // 优先级：1低，2中，3高
	Category          string     `json:"category"`          // 问题分类
	Title             string     `json:"title"`             // 工单标题
	Description       string     `json:"description"`       // 问题描述
	AttachUrls        string     `json:"attachUrls"`        // 附件URL列表
	Status            *int64     `json:"status"`            // 工单状态：1待受理，2处理中，3已完结，nil 不变更
	AssignedClerkID   *int64     `json:"assignedClerkId"`   // 分配客服ID
	AssignedClerkName string     `json:"assignedClerkName"` // 分配客服姓名
	ResolvedAt        *time.Time `json:"resolvedAt"`        // 解决时间，流转到已完结时为空则由后端补当前时间
	ClosedAt          *time.Time `json:"closedAt"`          // 关闭时间，流转到已完结时为空则由后端补当前时间
	Remark            string     `json:"remark"`            // 备注
}

// CsTicketIdReq 客服工单主键ID路径入参，用于 /api/cs/ticket/:id 类接口
type CsTicketIdReq struct {
	ID uint `json:"id" form:"id" uri:"id"` // 主键ID
}

func (r *CsTicketIdReq) Uint() uint {
	return r.ID
}

// CsTicketStatusReq 修改客服工单状态入参
// 仅变更工单状态及其关联的处理客服、解决/关闭时间、备注；工单号与其它业务字段不在此接口修改
// 状态流转：待受理(1) -> 处理中(2) -> 已完结(3)，已完结为终态
type CsTicketStatusReq struct {
	ID                uint       `json:"id"`                // 主键ID
	Status            int64      `json:"status"`            // 目标状态：1待受理，2处理中，3已完结
	AssignedClerkID   *int64     `json:"assignedClerkId"`   // 分配客服ID，流转到处理中时必填
	AssignedClerkName string     `json:"assignedClerkName"` // 分配客服姓名
	ResolvedAt        *time.Time `json:"resolvedAt"`        // 解决时间，流转到已完结时为空则由后端补当前时间
	ClosedAt          *time.Time `json:"closedAt"`          // 关闭时间，流转到已完结时为空则由后端补当前时间
	Remark            string     `json:"remark"`            // 备注
}

// ==========================================================
// cs_sessions 客服会话
// ==========================================================

// CsSessionSearch 客服会话分页查询入参
type CsSessionSearch struct {
	commonReq.PageInfo
	SessionNo string  `json:"sessionNo" form:"sessionNo"` // 会话编号，精确匹配
	Keyword   string  `json:"keyword" form:"keyword"`     // 综合关键词，会话编号/用户姓名/客服姓名 模糊匹配
	TicketID  *int64  `json:"ticketId" form:"ticketId"`   // 关联工单ID
	UserID    *int64  `json:"userId" form:"userId"`       // 用户ID
	UserPhone string  `json:"userPhone" form:"userPhone"` // 用户手机号，精确匹配
	ClerkID   *int64  `json:"clerkId" form:"clerkId"`     // 客服ID
	Status    *int64  `json:"status" form:"status"`       // 会话状态：1进行中，2已结束
}

// CsSessionCreateReq 新增客服会话入参
// 会话编号由后端 global.Snow 雪花ID自动生成，前端不传；新建会话固定为进行中
// TODO(ride-hailing)：会话如需展示订单/车辆信息，调用 ride-hailing 微服务「订单查询接口」获取
type CsSessionCreateReq struct {
	TicketID  *int64     `json:"ticketId"`  // 关联工单ID，无工单可不传
	UserID    int64      `json:"userId"`    // 用户ID
	UserName  string     `json:"userName"`  // 用户姓名
	UserPhone string     `json:"userPhone"` // 用户手机号
	ClerkID   *int64     `json:"clerkId"`   // 客服ID
	ClerkName string     `json:"clerkName"` // 客服姓名
	StartedAt *time.Time `json:"startedAt"` // 开始时间，不传取服务器当前时间
	Remark    string     `json:"remark"`    // 备注
}

// CsSessionUpdateReq 修改客服会话入参，会话编号不允许修改
// Status 使用指针：传 1 进行中 / 2 已结束，nil 表示不变更会话状态
type CsSessionUpdateReq struct {
	ID        uint       `json:"id"`        // 主键ID
	TicketID  *int64     `json:"ticketId"`  // 关联工单ID
	UserID    int64      `json:"userId"`    // 用户ID
	UserName  string     `json:"userName"`  // 用户姓名
	UserPhone string     `json:"userPhone"` // 用户手机号
	ClerkID   *int64     `json:"clerkId"`   // 客服ID
	ClerkName string     `json:"clerkName"` // 客服姓名
	Status    *int64     `json:"status"`    // 会话状态：1进行中，2已结束，nil 不变更
	StartedAt *time.Time `json:"startedAt"` // 开始时间
	EndedAt   *time.Time `json:"endedAt"`   // 结束时间，流转到已结束时为空则由后端补当前时间
	Remark    string     `json:"remark"`    // 备注
}

// ==========================================================
// cs_messages 客服会话消息
// ==========================================================

// CsMessageSearch 客服会话消息分页查询入参，传 SessionID 时按会话历史消息语义查询（消息序号正序）
type CsMessageSearch struct {
	commonReq.PageInfo
	SessionID  *int64 `json:"sessionId" form:"sessionId"`   // 归属会话ID
	SenderType *int64 `json:"senderType" form:"senderType"` // 发送方类型：1乘客，2客服，3系统
	MsgType    *int64 `json:"msgType" form:"msgType"`       // 消息类型：1文本，2图片，3文件，4系统通知
	IsRead     *int64 `json:"isRead" form:"isRead"`         // 已读状态：0未读，1已读
	Keyword    string `json:"keyword" form:"keyword"`       // 综合关键词，消息文本内容 模糊匹配
}

// CsMessageCreateReq 保存会话消息入参
// 消息序号不传时后端按当前会话最大序号+1自动生成，消息发送时间不传时取服务器当前时间
type CsMessageCreateReq struct {
	SessionID  int64      `json:"sessionId"`  // 归属会话ID
	SenderType int64      `json:"senderType"` // 发送方类型：1乘客，2客服，3系统
	MsgType    int64      `json:"msgType"`    // 消息类型：1文本，2图片，3文件，4系统通知
	Content    string     `json:"content"`    // 消息文本内容
	AttachUrl  string     `json:"attachUrl"`  // 附件地址
	Seq        int64      `json:"seq"`        // 消息序号，不传后端自动生成
	IsRead     int64      `json:"isRead"`     // 已读状态：0未读，1已读
	SendTime   *time.Time `json:"sendTime"`   // 消息发送时间，不传取服务器当前时间
}

// CsMessageUpdateReq 修改会话消息入参
type CsMessageUpdateReq struct {
	ID         uint       `json:"id"`         // 主键ID
	SenderType int64      `json:"senderType"` // 发送方类型：1乘客，2客服，3系统
	MsgType    int64      `json:"msgType"`    // 消息类型：1文本，2图片，3文件，4系统通知
	Content    string     `json:"content"`    // 消息文本内容
	AttachUrl  string     `json:"attachUrl"`  // 附件地址
	Seq        int64      `json:"seq"`        // 消息序号
	IsRead     int64      `json:"isRead"`     // 已读状态：0未读，1已读
	SendTime   *time.Time `json:"sendTime"`   // 消息发送时间
}

// ==========================================================
// cs_faqs 常见问题
// ==========================================================

// CsFaqSearch 常见问题分页查询入参
type CsFaqSearch struct {
	commonReq.PageInfo
	Question string  `json:"question" form:"question"` // 问题标题，模糊匹配
	Keyword  string  `json:"keyword" form:"keyword"`   // 综合关键词，问题标题/回答内容/标签 模糊匹配
	Category string  `json:"category" form:"category"` // 问题分类，模糊匹配
	Tags     string  `json:"tags" form:"tags"`         // 标签，模糊匹配
	Status   *int64  `json:"status" form:"status"`     // 状态：0禁用，1启用
}

// CsFaqCreateReq 新增常见问题入参
type CsFaqCreateReq struct {
	Question string  `json:"question"` // 问题标题
	Answer   string  `json:"answer"`   // 回答内容
	Category string  `json:"category"` // 问题分类
	Tags     string  `json:"tags"`     // 标签
	Sort     int64   `json:"sort"`     // 排序号，不传按库默认0
	Status   *int64  `json:"status"`   // 状态：0禁用，1启用；nil 表示按库默认1启用
	Remark   string  `json:"remark"`   // 备注
}

// CsFaqUpdateReq 修改常见问题入参，含上下架
// Status 使用指针：传 0 下架 / 1 上架，nil 表示不变更上架状态
type CsFaqUpdateReq struct {
	ID       uint   `json:"id"`       // 主键ID
	Question string `json:"question"` // 问题标题
	Answer   string `json:"answer"`   // 回答内容
	Category string `json:"category"` // 问题分类
	Tags     string `json:"tags"`     // 标签
	Sort     int64  `json:"sort"`     // 排序号
	Status   *int64 `json:"status"`   // 状态：0禁用，1启用，nil 不变更
	Remark   string `json:"remark"`   // 备注
}

// ==========================================================
// cs_faq_category 常见问题分类
// ==========================================================

// CsFaqCategorySearch 常见问题分类分页查询入参
type CsFaqCategorySearch struct {
	commonReq.PageInfo
	Title  string  `json:"title" form:"title"`   // 分类名称，模糊匹配
	Status *int64  `json:"status" form:"status"` // 状态：0禁用，1启用
}

// CsFaqCategoryCreateReq 新增常见问题分类入参
type CsFaqCategoryCreateReq struct {
	Title  string  `json:"title"`  // 分类名称
	Sort   int64   `json:"sort"`   // 排序号，不传按库默认0
	Status *int64  `json:"status"` // 状态：0禁用，1启用；nil 表示按库默认1启用
	Remark string  `json:"remark"` // 备注
}

// CsFaqCategoryUpdateReq 修改常见问题分类入参
// Status 使用指针：传 0 禁用 / 1 启用，nil 表示不变更状态
type CsFaqCategoryUpdateReq struct {
	ID     uint    `json:"id"`     // 主键ID
	Title  string  `json:"title"`  // 分类名称
	Sort   int64   `json:"sort"`   // 排序号
	Status *int64  `json:"status"` // 状态：0禁用，1启用，nil 不变更
	Remark string  `json:"remark"` // 备注
}

// ==========================================================
// cs_notice 消息通知
// ==========================================================

// CsNoticeSearch 消息通知分页查询入参
type CsNoticeSearch struct {
	commonReq.PageInfo
	UserID  *int64  `json:"userId" form:"userId"`   // 接收人用户ID
	Title   string  `json:"title" form:"title"`     // 消息标题，模糊匹配
	Keyword string  `json:"keyword" form:"keyword"` // 综合关键词，消息标题/消息内容 模糊匹配
	IsRead  *int64  `json:"isRead" form:"isRead"`   // 是否已读：0未读，1已读
	Type    *int64  `json:"type" form:"type"`       // 消息类型：1系统通知，2工单通知，3会话通知
}

// CsNoticeCreateReq 新增消息通知入参，已读状态固定为未读，由 /read 接口标记
type CsNoticeCreateReq struct {
	UserID  int64   `json:"userId"`  // 接收人用户ID
	Title   string  `json:"title"`   // 消息标题
	Content string  `json:"content"` // 消息内容
	Type    int64   `json:"type"`    // 消息类型：1系统通知，2工单通知，3会话通知
}

// CsNoticeUpdateReq 修改消息通知入参
// Type / IsRead 使用指针：传值表示变更，nil 表示不变更
type CsNoticeUpdateReq struct {
	ID      uint    `json:"id"`      // 主键ID
	UserID  *int64  `json:"userId"`  // 接收人用户ID，nil 不变更
	Title   string  `json:"title"`   // 消息标题
	Content string  `json:"content"` // 消息内容
	Type    *int64  `json:"type"`    // 消息类型：1系统通知，2工单通知，3会话通知，nil 不变更
	IsRead  *int64  `json:"isRead"`  // 是否已读：0未读，1已读，nil 不变更
}

// CsNoticeReadReq 标记消息通知已读入参
type CsNoticeReadReq struct {
	ID uint `json:"id"` // 主键ID
}
