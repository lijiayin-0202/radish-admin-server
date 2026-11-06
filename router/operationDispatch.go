package router

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// OperationDispatchRouter 运营调度管理路由
type OperationDispatchRouter struct{}

// InitOperationDispatchRouter 注册运营调度管理路由（前缀 /operationDispatch）
func (s *OperationDispatchRouter) InitOperationDispatchRouter(Router *gin.RouterGroup) {
	api := &v1.ApiGroupApp.OperationDispatchApi
	write := Router.Group("operationDispatch").Use(middleware.OperationRecord())
	read := Router.Group("operationDispatch")
	{
		write.POST("create", api.CreateDispatch)
		write.PUT("update", api.UpdateDispatch)
		write.DELETE("delete", api.DeleteDispatch)
		write.DELETE("deleteByIds", api.DeleteDispatchByIds)
		write.POST("cancel", api.CancelDispatch)
	}
	{
		read.GET("find", api.FindDispatch)
		read.GET("list", api.GetDispatchList)
	}
}
