package ride

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ride"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
)

// 后台通过 ride-service 的内部接口操作订单，不直接跨库写 robotaxi schema。
const (
	defaultRideServiceURL   = "http://127.0.0.1:8000"
	defaultRideServiceToken = "dev-internal-token"
)

type OrderService struct{}

func rideServiceURL() string {
	if value := strings.TrimRight(os.Getenv("RIDE_SERVICE_URL"), "/"); value != "" {
		return value
	}
	return defaultRideServiceURL
}

func rideServiceToken() string {
	if value := os.Getenv("RIDE_INTERNAL_TOKEN"); value != "" {
		return value
	}
	return defaultRideServiceToken
}

func callRideService(method, path string, params map[string]string, body any) (map[string]any, error) {
	headers := map[string]string{"X-Internal-Token": rideServiceToken()}
	response, err := request.HttpRequest(rideServiceURL()+path, method, headers, params, body)
	if err != nil {
		return nil, fmt.Errorf("调用订单服务失败：%w", err)
	}
	defer func() { _ = response.Body.Close() }()
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取订单服务响应失败：%w", err)
	}
	result := map[string]any{}
	if len(payload) > 0 {
		if err = json.Unmarshal(payload, &result); err != nil {
			return nil, fmt.Errorf("解析订单服务响应失败：%w", err)
		}
	}
	if response.StatusCode != http.StatusOK {
		message, _ := result["message"].(string)
		if message == "" {
			message = fmt.Sprintf("订单服务返回 %d", response.StatusCode)
		}
		return nil, fmt.Errorf("%s", message)
	}
	return result, nil
}

// ListOrders 查询订单分页列表。
func (s *OrderService) ListOrders(query ride.OrderListQuery) (ride.OrderListResult, error) {
	result := ride.OrderListResult{}
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
	if query.Status != "" {
		params["status"] = query.Status
	}
	if query.Keyword != "" {
		params["keyword"] = query.Keyword
	}
	payload, err := callRideService(http.MethodGet, "/api/v1/internal/orders", params, nil)
	if err != nil {
		return result, err
	}
	raw, _ := json.Marshal(payload)
	if err = json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("解析订单列表失败：%w", err)
	}
	if result.Items == nil {
		result.Items = []ride.OrderItem{}
	}
	return result, nil
}

// GetOrder 查询订单详情（含派单、位置、核验摘要）。
func (s *OrderService) GetOrder(orderID string) (map[string]any, error) {
	return callRideService(http.MethodGet, "/api/v1/internal/orders/"+orderID, nil, nil)
}

// Dispatch 后台派车：派单 + 车辆确认，可指定车辆。
func (s *OrderService) Dispatch(payload ride.DispatchRequest) (map[string]any, error) {
	body := map[string]any{"vehicleId": payload.VehicleID}
	return callRideService(http.MethodPost, "/api/v1/internal/orders/"+payload.OrderID+"/dispatch", nil, body)
}

// Advance 后台推进订单状态。
func (s *OrderService) Advance(payload ride.CommandRequest) (map[string]any, error) {
	body := map[string]any{"commandType": strings.ToUpper(payload.CommandType)}
	return callRideService(http.MethodPost, "/api/v1/internal/orders/"+payload.OrderID+"/commands", nil, body)
}
