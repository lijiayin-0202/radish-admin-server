package module3

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/module3"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/module3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationManagementAPI 提供给运营管理页面的模块 3 BFF 接口。
type OperationManagementAPI struct {
	service *service.OperationManagement
}

func NewOperationManagementAPI(service *service.OperationManagement) *OperationManagementAPI {
	return &OperationManagementAPI{service: service}
}

func (api *OperationManagementAPI) ListServiceAreas(c *gin.Context) {
	var input model.ServiceAreaListInput
	if !bindJSON(c, &input) {
		return
	}
	output, err := api.service.ListServiceAreas(c.Request.Context(), input)
	if !respondServiceError(c, err) {
		return
	}
	response.OkWithDetailed(response.PageResult{List: output.ServiceAreas, Total: output.Total, Page: input.Page, PageSize: input.PageSize}, "获取成功", c)
}

func (api *OperationManagementAPI) GetServiceArea(c *gin.Context) {
	var input model.GetServiceAreaInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := api.service.GetServiceArea(c.Request.Context(), input.ID)
	if respondServiceError(c, err) {
		response.OkWithDetailed(detail, "获取成功", c)
	}
}

func (api *OperationManagementAPI) CreateServiceArea(c *gin.Context) {
	var input model.CreateServiceAreaInput
	if !bindJSON(c, &input) {
		return
	}
	area, err := api.service.CreateServiceArea(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"serviceArea": area}, "创建成功", c)
	}
}

func (api *OperationManagementAPI) UpdateServiceArea(c *gin.Context) {
	var input model.UpdateServiceAreaInput
	if !bindJSON(c, &input) {
		return
	}
	area, err := api.service.UpdateServiceArea(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"serviceArea": area}, "更新成功", c)
	}
}

func (api *OperationManagementAPI) ChangeServiceAreaStatus(c *gin.Context) {
	var input model.ChangeServiceAreaStatusInput
	if !bindJSON(c, &input) {
		return
	}
	area, err := api.service.ChangeServiceAreaStatus(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"serviceArea": area}, "状态更新成功", c)
	}
}

func (api *OperationManagementAPI) ReplaceServiceAreaSchedules(c *gin.Context) {
	var input model.ReplaceSchedulesInput
	if !bindJSON(c, &input) {
		return
	}
	detail, err := api.service.ReplaceServiceAreaSchedules(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(detail, "运营时段保存成功", c)
	}
}

func (api *OperationManagementAPI) ListStations(c *gin.Context) {
	var input model.ListStationsInput
	if !bindJSON(c, &input) {
		return
	}
	stations, err := api.service.ListStations(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"list": stations}, "获取成功", c)
	}
}

func (api *OperationManagementAPI) CreateStation(c *gin.Context) {
	var input model.CreateStationInput
	if !bindJSON(c, &input) {
		return
	}
	station, err := api.service.CreateStation(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"station": station}, "创建成功", c)
	}
}

func (api *OperationManagementAPI) UpdateStation(c *gin.Context) {
	var input model.UpdateStationInput
	if !bindJSON(c, &input) {
		return
	}
	station, err := api.service.UpdateStation(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"station": station}, "更新成功", c)
	}
}

func (api *OperationManagementAPI) ChangeStationStatus(c *gin.Context) {
	var input model.ChangeStationStatusInput
	if !bindJSON(c, &input) {
		return
	}
	station, err := api.service.ChangeStationStatus(c.Request.Context(), input)
	if respondServiceError(c, err) {
		response.OkWithDetailed(gin.H{"station": station}, "状态更新成功", c)
	}
}

func bindJSON(c *gin.Context, input interface{}) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		response.FailWithMessage(err.Error(), c)
		return false
	}
	return true
}

func respondServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return true
	}
	if global.GVA_LOG != nil {
		global.GVA_LOG.Error("模块3服务调用失败", zap.Error(err))
	}
	response.FailWithMessage("模块3服务暂不可用", c)
	return false
}
