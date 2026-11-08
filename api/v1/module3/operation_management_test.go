package module3

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	module3Model "github.com/flipped-aurora/gin-vue-admin/server/model/module3"
	module3Service "github.com/flipped-aurora/gin-vue-admin/server/service/module3"
	"github.com/gin-gonic/gin"
)

type fakeOperationClient struct {
	module3Service.ClientAPI
}

func (fakeOperationClient) ListServiceAreas(context.Context, module3Model.ServiceAreaListInput) (module3Model.ServiceAreaListOutput, error) {
	return module3Model.ServiceAreaListOutput{
		ServiceAreas: []module3Model.ServiceArea{{ID: 1}},
		Total:        1,
	}, nil
}

func (fakeOperationClient) GetServiceArea(context.Context, int64) (module3Model.ServiceAreaDetail, error) {
	return module3Model.ServiceAreaDetail{
		ServiceArea: module3Model.ServiceArea{ID: 1, AreaName: "西湖区域"},
		Schedules:   []module3Model.ServiceAreaSchedule{{DayOfWeek: 1, StartMinute: 540, EndMinute: 1080, IsEnabled: true}},
	}, nil
}

func TestOperationManagementAPIListServiceAreas(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOperationManagementAPI(module3Service.NewOperationManagement(fakeOperationClient{}))
	router := gin.New()
	router.POST("/module3/service-area/list", handler.ListServiceAreas)

	request := httptest.NewRequest(http.MethodPost, "/module3/service-area/list", strings.NewReader(`{"page":1,"pageSize":10}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	want := `{"code":0,"data":{"list":[{"id":"1","areaCode":"","areaName":"","cityCode":"","cityName":"","boundaryGeojson":"","serviceNotice":"","status":0}],"total":1,"page":1,"pageSize":10},"msg":"获取成功"}`
	if strings.TrimSpace(recorder.Body.String()) != want {
		t.Fatalf("body = %s, want %s", recorder.Body.String(), want)
	}
}

func TestOperationManagementAPIGetServiceAreaReturnsSchedules(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOperationManagementAPI(module3Service.NewOperationManagement(fakeOperationClient{}))
	router := gin.New()
	router.POST("/module3/service-area/detail", handler.GetServiceArea)

	request := httptest.NewRequest(http.MethodPost, "/module3/service-area/detail", strings.NewReader(`{"id":"1"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"startMinute":540`) {
		t.Fatalf("response = %s, want service-area schedules", recorder.Body.String())
	}
}

type unavailableOperationClient struct {
	module3Service.ClientAPI
}

func (unavailableOperationClient) ListServiceAreas(context.Context, module3Model.ServiceAreaListInput) (module3Model.ServiceAreaListOutput, error) {
	return module3Model.ServiceAreaListOutput{}, errors.New("Get http://127.0.0.1:8000: connection refused")
}

func TestOperationManagementAPIRejectsInvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOperationManagementAPI(module3Service.NewOperationManagement(fakeOperationClient{}))
	router := gin.New()
	router.POST("/module3/service-area/list", handler.ListServiceAreas)

	request := httptest.NewRequest(http.MethodPost, "/module3/service-area/list", strings.NewReader(`{"page":0,"pageSize":10}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"code":7`) {
		t.Fatalf("response = %s, want validation failure", recorder.Body.String())
	}
}

func TestOperationManagementAPIRejectsInvalidServiceAreaStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOperationManagementAPI(module3Service.NewOperationManagement(fakeOperationClient{}))
	router := gin.New()
	router.POST("/module3/service-area/change-status", handler.ChangeServiceAreaStatus)

	request := httptest.NewRequest(http.MethodPost, "/module3/service-area/change-status", strings.NewReader(`{"id":"1","status":3}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"code":7`) {
		t.Fatalf("response = %s, want validation failure", recorder.Body.String())
	}
}

func TestOperationManagementAPIMasksUpstreamAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOperationManagementAPI(module3Service.NewOperationManagement(unavailableOperationClient{}))
	router := gin.New()
	router.POST("/module3/service-area/list", handler.ListServiceAreas)

	request := httptest.NewRequest(http.MethodPost, "/module3/service-area/list", strings.NewReader(`{"page":1,"pageSize":10}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "模块3服务暂不可用") {
		t.Fatalf("response = %s, want generic unavailable message", recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "127.0.0.1:8000") {
		t.Fatalf("response leaked internal address: %s", recorder.Body.String())
	}
}
