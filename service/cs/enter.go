package cs

// ServiceGroup 客服中心与客服运营端 service 组
// 注意：本模块所有 DB 操作统一使用 global.LUOBO_DB（luoboPgsql 业务库），严禁使用 global.GVA_DB
type ServiceGroup struct {
	CsFaqService
	CsTicketService
	CsSessionService
	CsMessageService
	CsFaqCategoryService
	CsNoticeService
}
