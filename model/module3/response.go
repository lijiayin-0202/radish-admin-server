package module3

import (
	"encoding/json"
	"fmt"
)

// ServiceAreaStatus 兼容微服务返回的 Proto 枚举文本与旧版数值状态。
type ServiceAreaStatus int32

func (status *ServiceAreaStatus) UnmarshalJSON(data []byte) error {
	var enumName string
	if err := json.Unmarshal(data, &enumName); err == nil {
		switch enumName {
		case "SERVICE_AREA_STATUS_TEMPORARILY_CLOSED":
			*status = 0
		case "SERVICE_AREA_STATUS_OPERATING":
			*status = 1
		case "SERVICE_AREA_STATUS_OFFLINE":
			*status = 2
		default:
			return fmt.Errorf("unsupported service area status: %q", enumName)
		}
		return nil
	}

	var numericStatus int32
	if err := json.Unmarshal(data, &numericStatus); err != nil {
		return fmt.Errorf("decode service area status: %w", err)
	}
	*status = ServiceAreaStatus(numericStatus)
	return nil
}

type ServiceArea struct {
	ID              int64             `json:"id,string"`
	AreaCode        string            `json:"areaCode"`
	AreaName        string            `json:"areaName"`
	CityCode        string            `json:"cityCode"`
	CityName        string            `json:"cityName"`
	BoundaryGeoJSON string            `json:"boundaryGeojson"`
	ServiceNotice   string            `json:"serviceNotice"`
	Status          ServiceAreaStatus `json:"status"`
}

type ServiceAreaSchedule struct {
	DayOfWeek   int32  `json:"dayOfWeek"`
	SpecialDate string `json:"specialDate"`
	StartMinute int32  `json:"startMinute"`
	EndMinute   int32  `json:"endMinute"`
	IsEnabled   bool   `json:"isEnabled"`
	Notice      string `json:"notice"`
}

type ServiceAreaDetail struct {
	ServiceArea ServiceArea           `json:"serviceArea"`
	Schedules   []ServiceAreaSchedule `json:"schedules"`
}

type ServiceAreaListOutput struct {
	ServiceAreas []ServiceArea `json:"serviceAreas"`
	Total        int64         `json:"total,string"`
}

type Station struct {
	ID               int64   `json:"id,string"`
	ServiceAreaID    int64   `json:"serviceAreaId,string"`
	StationCode      string  `json:"stationCode"`
	StationName      string  `json:"stationName"`
	AliasName        string  `json:"aliasName"`
	Address          string  `json:"address"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	StationType      int32   `json:"stationType"`
	ServiceDirection int32   `json:"serviceDirection"`
	DirectionHint    string  `json:"directionHint"`
	Status           int32   `json:"status"`
}
