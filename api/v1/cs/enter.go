package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

// ApiGroup 客服中心与客服运营端 api 组
type ApiGroup struct {
	CsFaqApi
	CsTicketApi
	CsSessionApi
	CsMessageApi
	CsFaqCategoryApi
	CsNoticeApi
}

var (
	csFaqService         = service.ServiceGroupApp.CsServiceGroup.CsFaqService
	csTicketService      = service.ServiceGroupApp.CsServiceGroup.CsTicketService
	csSessionService     = service.ServiceGroupApp.CsServiceGroup.CsSessionService
	csMessageService     = service.ServiceGroupApp.CsServiceGroup.CsMessageService
	csFaqCategoryService = service.ServiceGroupApp.CsServiceGroup.CsFaqCategoryService
	csNoticeService      = service.ServiceGroupApp.CsServiceGroup.CsNoticeService
)

// 本模块自用参数校验规则（定义在业务侧，不修改框架 utils/validator.go）
var (
	// 工单：优先级 1低，2中，3高；状态流转在 service 层校验（待受理 -> 处理中 -> 已完结）
	CsTicketCreateVerify = utils.Rules{
		"UserID":   {utils.NotEmpty()},
		"Priority": {utils.Ge("1"), utils.Le("3")},
		"Title":    {utils.NotEmpty()},
	}
	CsTicketUpdateVerify = utils.Rules{
		"ID":       {utils.NotEmpty()},
		"Priority": {utils.Ge("0"), utils.Le("3")}, // 0 表示不变更优先级
	}
	CsTicketStatusVerify = utils.Rules{
		"ID":     {utils.NotEmpty()},
		"Status": {utils.NotEmpty(), utils.Ge("1"), utils.Le("3")}, // 具体流转在 service 层校验（待受理 -> 处理中 -> 已完结）
	}

	// 会话：用户ID必填，会话编号由后端生成
	CsSessionCreateVerify = utils.Rules{
		"UserID":   {utils.NotEmpty()},
		"UserName": {utils.NotEmpty()},
	}
	CsSessionUpdateVerify = utils.Rules{"ID": {utils.NotEmpty()}}

	// 消息：发送方类型 1乘客/2客服/3系统，消息类型 1文本/2图片/3文件/4系统通知
	CsMessageCreateVerify = utils.Rules{
		"SessionID":  {utils.NotEmpty()},
		"SenderType": {utils.NotEmpty(), utils.Ge("1"), utils.Le("3")},
		"MsgType":    {utils.NotEmpty(), utils.Ge("1"), utils.Le("4")},
		"Content":    {utils.NotEmpty()},
	}
	CsMessageUpdateVerify = utils.Rules{"ID": {utils.NotEmpty()}}

	// FAQ：问题与答案必填，上下架状态在 service 层校验（0禁用，1启用）
	CsFaqCreateVerify = utils.Rules{
		"Question": {utils.NotEmpty()},
		"Answer":   {utils.NotEmpty()},
	}
	CsFaqUpdateVerify = utils.Rules{"ID": {utils.NotEmpty()}}

	// FAQ分类：分类名称必填，启用状态在 service 层校验（0禁用，1启用）
	CsFaqCategoryCreateVerify = utils.Rules{
		"Title": {utils.NotEmpty()},
	}
	CsFaqCategoryUpdateVerify = utils.Rules{"ID": {utils.NotEmpty()}}

	// 消息通知：接收人、标题、内容、消息类型必填；已读状态由 /read 接口标记
	CsNoticeCreateVerify = utils.Rules{
		"UserID":  {utils.NotEmpty()},
		"Title":   {utils.NotEmpty()},
		"Content": {utils.NotEmpty()},
		"Type":    {utils.NotEmpty(), utils.Ge("1"), utils.Le("3")}, // 1系统通知，2工单通知，3会话通知
	}
	CsNoticeUpdateVerify = utils.Rules{"ID": {utils.NotEmpty()}}
	CsNoticeReadVerify   = utils.Rules{"ID": {utils.NotEmpty()}}
)
