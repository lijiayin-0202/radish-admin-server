package ride

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ride"
)

type PassengerService struct{}

// ListPassengers 查询乘客分页列表，数据来自 ride-service 内部接口。
func (s *PassengerService) ListPassengers(query ride.PassengerListQuery) (ride.PassengerListResult, error) {
	result := ride.PassengerListResult{}
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
	if query.RealNameStatus != "" {
		params["realNameStatus"] = query.RealNameStatus
	}
	if query.AccountStatus != "" {
		params["accountStatus"] = query.AccountStatus
	}
	payload, err := callRideService(http.MethodGet, "/api/v1/internal/passengers", params, nil)
	if err != nil {
		return result, err
	}
	raw, _ := json.Marshal(payload)
	if err = json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("解析乘客列表失败：%w", err)
	}
	if result.Items == nil {
		result.Items = []ride.PassengerItem{}
	}
	return result, nil
}

// GetPassenger 查询乘客详情。
func (s *PassengerService) GetPassenger(passengerID string) (map[string]any, error) {
	if passengerID == "" {
		return nil, fmt.Errorf("passengerId 不能为空")
	}
	return callRideService(http.MethodGet, "/api/v1/internal/passengers/"+passengerID, nil, nil)
}
