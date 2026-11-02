package router

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// ServiceAreaRouter 运营区域路由
type ServiceAreaRouter struct{}

// InitServiceAreaRouter 注册运营区域路由（前缀 /serviceArea）
func (s *ServiceAreaRouter) InitServiceAreaRouter(Router *gin.RouterGroup) {
	api := &v1.ApiGroupApp.ServiceAreaApi
	write := Router.Group("serviceArea").Use(middleware.OperationRecord())
	read := Router.Group("serviceArea")
	{
		write.POST("create", api.CreateServiceArea)
		write.PUT("update", api.UpdateServiceArea)
		write.DELETE("delete", api.DeleteServiceArea)
		write.DELETE("deleteByIds", api.DeleteServiceAreasByIds)
	}
	{
		read.GET("find", api.FindServiceArea)
		read.GET("list", api.GetServiceAreaList)
	}
}
