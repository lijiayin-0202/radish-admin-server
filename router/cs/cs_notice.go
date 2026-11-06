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
		noticeRouter.POST("", csNoticeApi.CreateCsNotice)   // 新增站内公告
		noticeRouter.PUT("", csNoticeApi.UpdateCsNotice)    // 修改站内公告
		noticeRouter.DELETE("", csNoticeApi.DeleteCsNotice) // 删除站内公告
	}
	{
		noticeRouterWithoutRecord.GET("", csNoticeApi.GetCsNotice)         // 获取站内公告详情
		noticeRouterWithoutRecord.GET("list", csNoticeApi.GetCsNoticeList) // 分页获取站内公告列表
	}
}
