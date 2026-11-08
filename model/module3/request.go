package module3

// ServiceAreaListInput 是后台区域列表查询参数。
type ServiceAreaListInput struct {
	Page     int    `json:"page" binding:"required,min=1"`
	PageSize int    `json:"pageSize" binding:"required,min=1,max=100"`
	CityCode string `json:"cityCode"`
	Status   *int32 `json:"status" binding:"omitempty,oneof=0 1 2"`
}

// GetServiceAreaInput 只接收区域 ID，供编辑运营时段前读取现有规则。
type GetServiceAreaInput struct {
	ID int64 `json:"id,string" binding:"required"`
}

type CreateServiceAreaInput struct {
	AreaCode        string `json:"areaCode" binding:"required,max=64"`
	AreaName        string `json:"areaName" binding:"required,max=128"`
	CityCode        string `json:"cityCode" binding:"required,max=64"`
	CityName        string `json:"cityName" binding:"required,max=128"`
	BoundaryGeoJSON string `json:"boundaryGeojson" binding:"required"`
	ServiceNotice   string `json:"serviceNotice"`
	Status          int32  `json:"status" binding:"oneof=0 1 2"`
}

type UpdateServiceAreaInput struct {
	ID              int64  `json:"id,string" binding:"required"`
	AreaName        string `json:"areaName" binding:"required,max=128"`
	CityCode        string `json:"cityCode" binding:"required,max=64"`
	CityName        string `json:"cityName" binding:"required,max=128"`
	BoundaryGeoJSON string `json:"boundaryGeojson" binding:"required"`
	ServiceNotice   string `json:"serviceNotice"`
}

type ChangeServiceAreaStatusInput struct {
	ID     int64 `json:"id,string" binding:"required"`
	Status int32 `json:"status" binding:"oneof=0 1 2"`
}

type ReplaceSchedulesInput struct {
	ID        int64                 `json:"id,string" binding:"required"`
	Schedules []ServiceAreaSchedule `json:"schedules"`
}

type ListStationsInput struct {
	ServiceAreaID int64 `json:"serviceAreaId,string" binding:"required"`
}

type CreateStationInput struct {
	ServiceAreaID    int64   `json:"serviceAreaId,string" binding:"required"`
	StationCode      string  `json:"stationCode" binding:"required,max=64"`
	StationName      string  `json:"stationName" binding:"required,max=128"`
	AliasName        string  `json:"aliasName"`
	Address          string  `json:"address" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	StationType      int32   `json:"stationType"`
	ServiceDirection int32   `json:"serviceDirection"`
	DirectionHint    string  `json:"directionHint"`
}

type UpdateStationInput struct {
	ID               int64   `json:"id,string" binding:"required"`
	StationName      string  `json:"stationName" binding:"required,max=128"`
	AliasName        string  `json:"aliasName"`
	Address          string  `json:"address" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	StationType      int32   `json:"stationType"`
	ServiceDirection int32   `json:"serviceDirection"`
	DirectionHint    string  `json:"directionHint"`
}

type ChangeStationStatusInput struct {
	ID     int64 `json:"id,string" binding:"required"`
	Status int32 `json:"status" binding:"oneof=0 1"`
}
