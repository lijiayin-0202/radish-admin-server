package module3

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1/module3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/module3"
	"github.com/gin-gonic/gin"
)

// OperationManagementRouter 将模块 3 运营管理接口挂载到后台私有路由组。
type OperationManagementRouter struct{}

func (router *OperationManagementRouter) InitOperationManagementRouter(group *gin.RouterGroup) {
	client := service.NewClient(global.GVA_CONFIG.RideHailing)
	handler := api.NewOperationManagementAPI(service.NewOperationManagement(client))
	routes := group.Group("module3")
	routes.POST("service-area/list", handler.ListServiceAreas)
	routes.POST("service-area/detail", handler.GetServiceArea)
	routes.POST("service-area/create", handler.CreateServiceArea)
	routes.POST("service-area/update", handler.UpdateServiceArea)
	routes.POST("service-area/change-status", handler.ChangeServiceAreaStatus)
	routes.POST("service-area/schedules", handler.ReplaceServiceAreaSchedules)
	routes.POST("station/list", handler.ListStations)
	routes.POST("station/create", handler.CreateStation)
	routes.POST("station/update", handler.UpdateStation)
	routes.POST("station/change-status", handler.ChangeStationStatus)
}
