package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	// cs 客服模块：统一挂在 /api/cs 前缀下，继承 privateGroup 的 JWT + Casbin 鉴权
	csGroup := privateGroup.Group("/api/cs")
	csRouter := router.RouterGroupApp.Cs
	csRouter.InitCsTicketRouter(csGroup)
	csRouter.InitCsSessionRouter(csGroup)
	csRouter.InitCsMessageRouter(csGroup)
	csRouter.InitCsFaqRouter(csGroup)
	csRouter.InitCsFaqCategoryRouter(csGroup)
	csRouter.InitCsNoticeRouter(csGroup)

	holder(publicGroup, privateGroup)
}
