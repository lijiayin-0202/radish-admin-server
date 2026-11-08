package router

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

// StationRouter 站点管理路由
type StationRouter struct{}

// InitStationRouter 注册站点管理路由（前缀 /station）
func (s *StationRouter) InitStationRouter(Router *gin.RouterGroup) {
	api := &v1.ApiGroupApp.StationApi
	write := Router.Group("station").Use(middleware.OperationRecord())
	read := Router.Group("station")
	{
		write.POST("create", api.CreateStation)
		write.PUT("update", api.UpdateStation)
		write.DELETE("delete", api.DeleteStation)
		write.DELETE("deleteByIds", api.DeleteStationsByIds)
	}
	{
		read.GET("find", api.FindStation)
		read.GET("list", api.GetStationList)
	}
}
