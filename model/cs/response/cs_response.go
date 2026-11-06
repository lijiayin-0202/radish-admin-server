package response

import (
	"time"

	csModel "github.com/flipped-aurora/gin-vue-admin/server/model/cs"
)

// CsTicketDto 工单返回 DTO
// 通过对外层 CreatedAt / UpdatedAt 重声明 json 标签，将时间字段输出为驼峰 createdAt / updatedAt
// 这样只影响 cs 工单接口，不改动全局 GVA_MODEL 的 json 标签，避免波及其他系统页面
type CsTicketDto struct {
	csModel.CsTicket
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsTicketDto 将工单实体转换为返回 DTO，时间字段驼峰化
func ToCsTicketDto(e csModel.CsTicket) CsTicketDto {
	return CsTicketDto{
		CsTicket:  e,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToCsTicketDtoList 批量转换工单实体为 DTO
func ToCsTicketDtoList(entities []csModel.CsTicket) []CsTicketDto {
	list := make([]CsTicketDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsTicketDto(e))
	}
	return list
}

// CsFaqDto FAQ问答返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsFaqDto struct {
	csModel.CsFaq
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsFaqDto 将 FAQ 实体转换为返回 DTO，时间字段驼峰化
func ToCsFaqDto(e csModel.CsFaq) CsFaqDto {
	return CsFaqDto{
		CsFaq:     e,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToCsFaqDtoList 批量转换 FAQ 实体为 DTO
func ToCsFaqDtoList(entities []csModel.CsFaq) []CsFaqDto {
	list := make([]CsFaqDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsFaqDto(e))
	}
	return list
}

// CsFaqCategoryDto 常见问题分类返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsFaqCategoryDto struct {
	csModel.CsFaqCategory
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsFaqCategoryDto 将分类实体转换为返回 DTO，时间字段驼峰化
func ToCsFaqCategoryDto(e csModel.CsFaqCategory) CsFaqCategoryDto {
	return CsFaqCategoryDto{
		CsFaqCategory: e,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// ToCsFaqCategoryDtoList 批量转换分类实体为 DTO
func ToCsFaqCategoryDtoList(entities []csModel.CsFaqCategory) []CsFaqCategoryDto {
	list := make([]CsFaqCategoryDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsFaqCategoryDto(e))
	}
	return list
}

// CsSessionDto 客服会话返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsSessionDto struct {
	csModel.CsSession
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsSessionDto 将会话实体转换为返回 DTO，时间字段驼峰化
func ToCsSessionDto(e csModel.CsSession) CsSessionDto {
	return CsSessionDto{
		CsSession: e,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToCsSessionDtoList 批量转换会话实体为 DTO
func ToCsSessionDtoList(entities []csModel.CsSession) []CsSessionDto {
	list := make([]CsSessionDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsSessionDto(e))
	}
	return list
}

// CsMessageDto 客服会话消息返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsMessageDto struct {
	csModel.CsMessage
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsMessageDto 将消息实体转换为返回 DTO，时间字段驼峰化
func ToCsMessageDto(e csModel.CsMessage) CsMessageDto {
	return CsMessageDto{
		CsMessage: e,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToCsMessageDtoList 批量转换消息实体为 DTO
func ToCsMessageDtoList(entities []csModel.CsMessage) []CsMessageDto {
	list := make([]CsMessageDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsMessageDto(e))
	}
	return list
}

// CsTicketResponse 客服工单详情返回
type CsTicketResponse struct {
	Ticket CsTicketDto `json:"ticket"`
}

// CsSessionResponse 客服会话详情返回
type CsSessionResponse struct {
	Session CsSessionDto `json:"session"`
}

// CsSessionCreateResponse 创建客服会话返回，会话编号由后端生成，前端据此发起后续会话消息
type CsSessionCreateResponse struct {
	SessionID uint   `json:"sessionId"`
	SessionNo string `json:"sessionNo"`
}

// CsMessageResponse 客服会话消息详情返回
type CsMessageResponse struct {
	Message CsMessageDto `json:"message"`
}

// CsFaqResponse 常见问题详情返回
type CsFaqResponse struct {
	Faq CsFaqDto `json:"faq"`
}

// CsFaqCategoryResponse 常见问题分类详情返回
type CsFaqCategoryResponse struct {
	Category CsFaqCategoryDto `json:"category"`
}

// CsNotificationDto 用户通知返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsNotificationDto struct {
	csModel.CsNotification
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsNotificationDto 将用户通知实体转换为返回 DTO，时间字段驼峰化
func ToCsNotificationDto(e csModel.CsNotification) CsNotificationDto {
	return CsNotificationDto{
		CsNotification: e,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

// ToCsNotificationDtoList 批量转换用户通知实体为 DTO
func ToCsNotificationDtoList(entities []csModel.CsNotification) []CsNotificationDto {
	list := make([]CsNotificationDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsNotificationDto(e))
	}
	return list
}

// CsNotificationResponse 用户通知详情返回
type CsNotificationResponse struct {
	Notification CsNotificationDto `json:"notification"`
}

// CsNoticeDto 站内公告返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsNoticeDto struct {
	csModel.CsNotice
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsNoticeDto 将站内公告实体转换为返回 DTO，时间字段驼峰化
func ToCsNoticeDto(e csModel.CsNotice) CsNoticeDto {
	return CsNoticeDto{
		CsNotice:  e,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// ToCsNoticeDtoList 批量转换站内公告实体为 DTO
func ToCsNoticeDtoList(entities []csModel.CsNotice) []CsNoticeDto {
	list := make([]CsNoticeDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsNoticeDto(e))
	}
	return list
}

// CsNoticeResponse 站内公告详情返回
type CsNoticeResponse struct {
	Notice CsNoticeDto `json:"notice"`
}

// CsTicketLogDto 工单流水返回 DTO，时间字段输出驼峰 createdAt/updatedAt，不改动全局 GVA_MODEL
type CsTicketLogDto struct {
	csModel.CsTicketLog
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToCsTicketLogDto 将工单流水实体转换为返回 DTO，时间字段驼峰化
func ToCsTicketLogDto(e csModel.CsTicketLog) CsTicketLogDto {
	return CsTicketLogDto{
		CsTicketLog: e,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToCsTicketLogDtoList 批量转换工单流水实体为 DTO
func ToCsTicketLogDtoList(entities []csModel.CsTicketLog) []CsTicketLogDto {
	list := make([]CsTicketLogDto, 0, len(entities))
	for _, e := range entities {
		list = append(list, ToCsTicketLogDto(e))
	}
	return list
}

// CsTicketLogResponse 工单流水详情返回
type CsTicketLogResponse struct {
	Log CsTicketLogDto `json:"log"`
}
