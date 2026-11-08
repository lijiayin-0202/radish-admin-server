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

	// 平台配置、权限与统计大屏
	router.RouterGroupApp.Platform.InitPlatformRouter(privateGroup)

	// 出行运营（模块1-4 后台视角）
	router.RouterGroupApp.Trip.InitTripRouter(privateGroup)

	// 运营调度管理
	router.RouterGroupApp.OperationDispatchRouter.InitOperationDispatchRouter(privateGroup)

	// cs 客服模块：InitCsRouter 内部自行 Group("cs")，因此这里直接传 privateGroup。
	// 前端 VITE_BASE_API=/api，Vite 代理会 rewrite 掉 /api 后再转发，
	// 后端路由不能再带 /api，否则 /cs/... 请求会 404。
	router.RouterGroupApp.Cs.InitCsRouter(privateGroup)

	holder(publicGroup, privateGroup)
}
