package cs

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CsNoticeRouter struct{}

func (r *CsNoticeRouter) InitCsNoticeRouter(Router *gin.RouterGroup) {
	noticeRouter := Router.Group("notice").Use(middleware.OperationRecord())
	noticeRouterWithoutRecord := Router.Group("notice")
	{
		noticeRouter.POST("", csNoticeApi.CreateCsNotice)    // 新增消息通知
		noticeRouter.PUT("", csNoticeApi.UpdateCsNotice)     // 修改消息通知
		noticeRouter.PUT("read", csNoticeApi.ReadCsNotice)   // 标记消息通知已读
		noticeRouter.DELETE("", csNoticeApi.DeleteCsNotice)  // 删除消息通知
	}
	{
		noticeRouterWithoutRecord.GET("", csNoticeApi.GetCsNotice)         // 获取消息通知详情
		noticeRouterWithoutRecord.GET("list", csNoticeApi.GetCsNoticeList) // 分页获取消息通知列表
	}
}
