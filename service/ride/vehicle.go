package ride

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ride"
)

type VehicleService struct{}

// ListVehicles 查询车辆分页列表，数据来自 ride-service 内部接口。
func (s *VehicleService) ListVehicles(query ride.VehicleListQuery) (ride.VehicleListResult, error) {
	result := ride.VehicleListResult{}
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	params := map[string]string{
		"page":     strconv.Itoa(page),
		"pageSize": strconv.Itoa(pageSize),
	}
	if query.Keyword != "" {
		params["keyword"] = query.Keyword
	}
	if query.OperationStatus != "" {
		params["operationStatus"] = query.OperationStatus
	}
	if query.OnlineStatus != "" {
		params["onlineStatus"] = query.OnlineStatus
	}
	payload, err := callRideService(http.MethodGet, "/api/v1/internal/vehicles", params, nil)
	if err != nil {
		return result, err
	}
	raw, _ := json.Marshal(payload)
	if err = json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("解析车辆列表失败：%w", err)
	}
	if result.Items == nil {
		result.Items = []ride.VehicleItem{}
	}
	return result, nil
}

// GetVehicle 查询车辆详情。
func (s *VehicleService) GetVehicle(vehicleID string) (map[string]any, error) {
	if vehicleID == "" {
		return nil, fmt.Errorf("vehicleId 不能为空")
	}
	return callRideService(http.MethodGet, "/api/v1/internal/vehicles/"+vehicleID, nil, nil)
}
