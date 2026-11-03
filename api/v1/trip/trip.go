package trip

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
	tripReq "github.com/flipped-aurora/gin-vue-admin/server/model/trip/request"
	tripService "github.com/flipped-aurora/gin-vue-admin/server/service/trip"
	"github.com/gin-gonic/gin"
)

type TripApi struct{}

func strToUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

// ==================== 乘客 ====================

func (a *TripApi) GetPassengerList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetPassengerList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreatePassenger(c *gin.Context) {
	var req tripModel.Passenger
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreatePassenger(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *TripApi) UpdatePassenger(c *gin.Context) {
	var req tripModel.Passenger
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdatePassenger(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *TripApi) DeletePassenger(c *gin.Context) {
	if err := tripService.TripServiceApp.DeletePassenger(strToUint(c.Param("id"))); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== 订单 ====================

func (a *TripApi) GetOrderList(c *gin.Context) {
	var info tripReq.OrderSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetOrderList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) GetOrderStatusLogs(c *gin.Context) {
	var q tripReq.OrderLogQuery
	_ = c.ShouldBindQuery(&q)
	list, err := tripService.TripServiceApp.GetOrderStatusLogs(q.OrderID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// ==================== 运营区域 ====================

func (a *TripApi) GetAreaList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetAreaList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreateArea(c *gin.Context) {
	var req tripModel.ServiceArea
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreateArea(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *TripApi) UpdateArea(c *gin.Context) {
	var req tripModel.ServiceArea
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateArea(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *TripApi) DeleteArea(c *gin.Context) {
	if err := tripService.TripServiceApp.DeleteArea(strToUint(c.Param("id"))); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== 站点 ====================

func (a *TripApi) GetStationList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetStationList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreateStation(c *gin.Context) {
	var req tripModel.Station
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreateStation(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *TripApi) UpdateStation(c *gin.Context) {
	var req tripModel.Station
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateStation(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *TripApi) DeleteStation(c *gin.Context) {
	if err := tripService.TripServiceApp.DeleteStation(strToUint(c.Param("id"))); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== 车辆 ====================

func (a *TripApi) GetVehicleList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetVehicleList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreateVehicle(c *gin.Context) {
	var req tripModel.Vehicle
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreateVehicle(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *TripApi) UpdateVehicle(c *gin.Context) {
	var req tripModel.Vehicle
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateVehicle(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *TripApi) DeleteVehicle(c *gin.Context) {
	if err := tripService.TripServiceApp.DeleteVehicle(strToUint(c.Param("id"))); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// ==================== 派单记录 ====================

func (a *TripApi) GetDispatchList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetDispatchList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// ==================== 车辆指令 ====================

func (a *TripApi) GetCommandList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetCommandList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreateCommand(c *gin.Context) {
	var req tripModel.VehicleCommand
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreateCommand(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("下发成功", c)
}

func (a *TripApi) UpdateCommand(c *gin.Context) {
	var req tripModel.VehicleCommand
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateCommand(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ==================== 车辆故障 ====================

func (a *TripApi) GetFaultList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetFaultList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) CreateFault(c *gin.Context) {
	var req tripModel.VehicleFault
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.CreateFault(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("登记成功", c)
}

func (a *TripApi) UpdateFault(c *gin.Context) {
	var req tripModel.VehicleFault
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateFault(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// ==================== 上车核验 ====================

func (a *TripApi) GetBoardingList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetBoardingList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

// ==================== 安全事件 ====================

func (a *TripApi) GetIncidentList(c *gin.Context) {
	var info tripReq.TripSearch
	_ = c.ShouldBindQuery(&info)
	list, total, err := tripService.TripServiceApp.GetIncidentList(info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: info.Page, PageSize: info.PageSize}, "获取成功", c)
}

func (a *TripApi) UpdateIncident(c *gin.Context) {
	var req tripModel.SafetyIncident
	_ = c.ShouldBindJSON(&req)
	if err := tripService.TripServiceApp.UpdateIncident(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
