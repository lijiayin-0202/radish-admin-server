package ride

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ride"
)

type DispatchService struct{}

// ListDispatchRecords 查询派单记录分页列表，数据来自 ride-service 内部接口。
func (s *DispatchService) ListDispatchRecords(query ride.DispatchRecordListQuery) (ride.DispatchRecordListResult, error) {
	result := ride.DispatchRecordListResult{}
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
	if query.Status != "" {
		params["status"] = query.Status
	}
	payload, err := callRideService(http.MethodGet, "/api/v1/internal/dispatch-records", params, nil)
	if err != nil {
		return result, err
	}
	raw, _ := json.Marshal(payload)
	if err = json.Unmarshal(raw, &result); err != nil {
		return result, fmt.Errorf("解析派单记录失败：%w", err)
	}
	if result.Items == nil {
		result.Items = []ride.DispatchRecordItem{}
	}
	return result, nil
}
