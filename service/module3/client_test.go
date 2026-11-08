package module3

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
)

func TestClientListServiceAreasUsesQueryParameters(t *testing.T) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", request.Method, http.MethodGet)
		}
		if request.URL.Path != "/v1/admin/service-areas" {
			t.Fatalf("path = %s, want /v1/admin/service-areas", request.URL.Path)
		}
		if got := request.URL.Query().Get("page"); got != "1" {
			t.Fatalf("page = %q, want 1", got)
		}
		if got := request.URL.Query().Get("page_size"); got != "10" {
			t.Fatalf("page_size = %q, want 10", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"serviceAreas":[],"total":"0"}`))
	}))
	defer server.Close()

	client := NewClient(config.RideHailing{BaseURL: server.URL, RequestTimeout: time.Second})
	areas, err := client.ListServiceAreas(context.Background(), ServiceAreaListInput{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("ListServiceAreas() error = %v", err)
	}
	if areas.ServiceAreas == nil {
		t.Fatal("ListServiceAreas() returned nil serviceAreas, want decoded empty list")
	}
	if areas.Total != 0 {
		t.Fatalf("ListServiceAreas() total = %d, want 0", areas.Total)
	}
}

func TestClientListServiceAreasDecodesProtoEnumStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"serviceAreas":[{"id":"101","areaName":"西湖区域","status":"SERVICE_AREA_STATUS_OPERATING"}],"total":"1"}`))
	}))
	defer server.Close()

	client := NewClient(config.RideHailing{BaseURL: server.URL, RequestTimeout: time.Second})
	areas, err := client.ListServiceAreas(context.Background(), ServiceAreaListInput{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("ListServiceAreas() error = %v", err)
	}
	if len(areas.ServiceAreas) != 1 {
		t.Fatalf("ListServiceAreas() count = %d, want 1", len(areas.ServiceAreas))
	}
	if areas.ServiceAreas[0].Status != 1 {
		t.Fatalf("ListServiceAreas() status = %d, want 1", areas.ServiceAreas[0].Status)
	}
}

func TestClientCreateServiceAreaSendsBodyAndReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("method = %s, want %s", request.Method, http.MethodPost)
		}
		if request.URL.Path != "/v1/admin/service-areas" {
			t.Fatalf("path = %s, want /v1/admin/service-areas", request.URL.Path)
		}
		var payload CreateServiceAreaInput
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if payload.AreaCode != "HZ-01" {
			t.Fatalf("areaCode = %q, want HZ-01", payload.AreaCode)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"serviceArea":{"id":"101","areaCode":"HZ-01","areaName":"西湖区域"}}`))
	}))
	defer server.Close()

	client := NewClient(config.RideHailing{BaseURL: server.URL, RequestTimeout: time.Second})
	area, err := client.CreateServiceArea(context.Background(), CreateServiceAreaInput{
		AreaCode: "HZ-01", AreaName: "西湖区域", CityCode: "330100", CityName: "杭州", BoundaryGeoJSON: "{}",
	})
	if err != nil {
		t.Fatalf("CreateServiceArea() error = %v", err)
	}
	if area.ID != 101 || area.AreaName != "西湖区域" {
		t.Fatalf("CreateServiceArea() = %#v, want decoded service area", area)
	}
}

func TestClientGetServiceAreaReturnsSchedules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", request.Method, http.MethodGet)
		}
		if request.URL.Path != "/v1/admin/service-areas/101" {
			t.Fatalf("path = %s, want /v1/admin/service-areas/101", request.URL.Path)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"serviceArea":{"id":"101","areaName":"西湖区域"},"schedules":[{"dayOfWeek":1,"startMinute":540,"endMinute":1080,"isEnabled":true}]}`))
	}))
	defer server.Close()

	client := NewClient(config.RideHailing{BaseURL: server.URL, RequestTimeout: time.Second})
	detail, err := client.GetServiceArea(context.Background(), 101)
	if err != nil {
		t.Fatalf("GetServiceArea() error = %v", err)
	}
	if detail.ServiceArea.ID != 101 || len(detail.Schedules) != 1 || detail.Schedules[0].StartMinute != 540 {
		t.Fatalf("GetServiceArea() = %#v, want decoded detail", detail)
	}
}

func TestClientReturnsGenericErrorWhenBaseURLIsEmpty(t *testing.T) {
	client := NewClient(config.RideHailing{})
	_, err := client.ListServiceAreas(context.Background(), ServiceAreaListInput{Page: 1, PageSize: 10})
	if !errors.Is(err, ErrRideHailingUnavailable) {
		t.Fatalf("ListServiceAreas() error = %v, want generic unavailable error", err)
	}
	if err.Error() != "模块3服务暂不可用" {
		t.Fatalf("error = %q, want no internal URL", err)
	}
}

func TestClientReturnsGenericErrorForUpstreamFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte("upstream is unavailable"))
	}))
	defer server.Close()

	client := NewClient(config.RideHailing{BaseURL: server.URL, RequestTimeout: time.Second})
	_, err := client.ListServiceAreas(context.Background(), ServiceAreaListInput{Page: 1, PageSize: 10})
	if !errors.Is(err, ErrRideHailingUnavailable) {
		t.Fatalf("ListServiceAreas() error = %v, want generic unavailable error", err)
	}
	if err.Error() != "模块3服务暂不可用" {
		t.Fatalf("error = %q, want no upstream detail", err)
	}
}
