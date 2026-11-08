package module3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/module3"
)

var ErrRideHailingUnavailable = errors.New("模块3服务暂不可用")

// ClientAPI 定义模块 3 微服务对后台管理端开放的 HTTP 调用。
type ClientAPI interface {
	ListServiceAreas(context.Context, ServiceAreaListInput) (ServiceAreaListOutput, error)
	GetServiceArea(context.Context, int64) (ServiceAreaDetail, error)
	CreateServiceArea(context.Context, CreateServiceAreaInput) (ServiceArea, error)
	UpdateServiceArea(context.Context, UpdateServiceAreaInput) (ServiceArea, error)
	ChangeServiceAreaStatus(context.Context, ChangeServiceAreaStatusInput) (ServiceArea, error)
	ReplaceServiceAreaSchedules(context.Context, ReplaceSchedulesInput) (ServiceAreaDetail, error)
	ListStations(context.Context, ListStationsInput) ([]Station, error)
	CreateStation(context.Context, CreateStationInput) (Station, error)
	UpdateStation(context.Context, UpdateStationInput) (Station, error)
	ChangeStationStatus(context.Context, ChangeStationStatusInput) (Station, error)
}

type Client struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

type ServiceAreaListInput = model.ServiceAreaListInput
type ServiceAreaListOutput = model.ServiceAreaListOutput
type CreateServiceAreaInput = model.CreateServiceAreaInput
type UpdateServiceAreaInput = model.UpdateServiceAreaInput
type ChangeServiceAreaStatusInput = model.ChangeServiceAreaStatusInput
type ReplaceSchedulesInput = model.ReplaceSchedulesInput
type ListStationsInput = model.ListStationsInput
type CreateStationInput = model.CreateStationInput
type UpdateStationInput = model.UpdateStationInput
type ChangeStationStatusInput = model.ChangeStationStatusInput
type ServiceArea = model.ServiceArea
type ServiceAreaDetail = model.ServiceAreaDetail
type Station = model.Station

func NewClient(configuration config.RideHailing) *Client {
	timeout := configuration.RequestTimeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &Client{baseURL: strings.TrimRight(configuration.BaseURL, "/"), timeout: timeout, httpClient: &http.Client{Timeout: timeout}}
}

func (client *Client) ListServiceAreas(ctx context.Context, input ServiceAreaListInput) (ServiceAreaListOutput, error) {
	query := url.Values{"page": {strconv.Itoa(input.Page)}, "page_size": {strconv.Itoa(input.PageSize)}}
	if input.CityCode != "" {
		query.Set("city_code", input.CityCode)
	}
	if input.Status != nil {
		query.Set("status", strconv.FormatInt(int64(*input.Status), 10))
	}
	var output ServiceAreaListOutput
	return output, client.do(ctx, http.MethodGet, "/v1/admin/service-areas", query, nil, &output)
}

// GetServiceArea 读取区域资料与当前完整运营时段，用于后台编辑时回填表单。
func (client *Client) GetServiceArea(ctx context.Context, id int64) (ServiceAreaDetail, error) {
	var output ServiceAreaDetail
	path := "/v1/admin/service-areas/" + strconv.FormatInt(id, 10)
	return output, client.do(ctx, http.MethodGet, path, nil, nil, &output)
}

func (client *Client) CreateServiceArea(ctx context.Context, input CreateServiceAreaInput) (ServiceArea, error) {
	var output struct {
		ServiceArea ServiceArea `json:"serviceArea"`
	}
	return output.ServiceArea, client.do(ctx, http.MethodPost, "/v1/admin/service-areas", nil, input, &output)
}

func (client *Client) UpdateServiceArea(ctx context.Context, input UpdateServiceAreaInput) (ServiceArea, error) {
	var output struct {
		ServiceArea ServiceArea `json:"serviceArea"`
	}
	return output.ServiceArea, client.do(ctx, http.MethodPatch, "/v1/admin/service-areas/"+strconv.FormatInt(input.ID, 10), nil, input, &output)
}

func (client *Client) ChangeServiceAreaStatus(ctx context.Context, input ChangeServiceAreaStatusInput) (ServiceArea, error) {
	var output struct {
		ServiceArea ServiceArea `json:"serviceArea"`
	}
	path := "/v1/admin/service-areas/" + strconv.FormatInt(input.ID, 10) + "/status"
	return output.ServiceArea, client.do(ctx, http.MethodPatch, path, nil, input, &output)
}

func (client *Client) ReplaceServiceAreaSchedules(ctx context.Context, input ReplaceSchedulesInput) (ServiceAreaDetail, error) {
	var output ServiceAreaDetail
	path := "/v1/admin/service-areas/" + strconv.FormatInt(input.ID, 10) + "/schedules"
	return output, client.do(ctx, http.MethodPut, path, nil, input, &output)
}

func (client *Client) ListStations(ctx context.Context, input ListStationsInput) ([]Station, error) {
	query := url.Values{"service_area_id": {strconv.FormatInt(input.ServiceAreaID, 10)}}
	var output struct {
		Stations []Station `json:"stations"`
	}
	return output.Stations, client.do(ctx, http.MethodGet, "/v1/admin/stations", query, nil, &output)
}

func (client *Client) CreateStation(ctx context.Context, input CreateStationInput) (Station, error) {
	var output struct {
		Station Station `json:"station"`
	}
	return output.Station, client.do(ctx, http.MethodPost, "/v1/admin/stations", nil, input, &output)
}

func (client *Client) UpdateStation(ctx context.Context, input UpdateStationInput) (Station, error) {
	var output struct {
		Station Station `json:"station"`
	}
	return output.Station, client.do(ctx, http.MethodPatch, "/v1/admin/stations/"+strconv.FormatInt(input.ID, 10), nil, input, &output)
}

func (client *Client) ChangeStationStatus(ctx context.Context, input ChangeStationStatusInput) (Station, error) {
	var output struct {
		Station Station `json:"station"`
	}
	path := "/v1/admin/stations/" + strconv.FormatInt(input.ID, 10) + "/status"
	return output.Station, client.do(ctx, http.MethodPatch, path, nil, input, &output)
}

func (client *Client) do(ctx context.Context, method, path string, query url.Values, payload any, output any) error {
	if client.baseURL == "" {
		return ErrRideHailingUnavailable
	}
	target, err := url.Parse(client.baseURL + path)
	if err != nil {
		return ErrRideHailingUnavailable
	}
	target.RawQuery = query.Encode()
	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		encoded, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return fmt.Errorf("%w: 请求编码失败", ErrRideHailingUnavailable)
		}
		body = bytes.NewReader(encoded)
	}
	requestContext, cancel := context.WithTimeout(ctx, client.timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestContext, method, target.String(), body)
	if err != nil {
		return ErrRideHailingUnavailable
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.httpClient.Do(request)
	if err != nil {
		return ErrRideHailingUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return ErrRideHailingUnavailable
	}
	if err = json.NewDecoder(response.Body).Decode(output); err != nil {
		return ErrRideHailingUnavailable
	}
	return nil
}
