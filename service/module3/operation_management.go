package module3

import "context"

// OperationManagement 只负责把 BFF 请求转发给模块 3 客户端，不包含持久化逻辑。
type OperationManagement struct {
	client ClientAPI
}

func NewOperationManagement(client ClientAPI) *OperationManagement {
	return &OperationManagement{client: client}
}

func (service *OperationManagement) ListServiceAreas(ctx context.Context, input ServiceAreaListInput) (ServiceAreaListOutput, error) {
	return service.client.ListServiceAreas(ctx, input)
}

func (service *OperationManagement) GetServiceArea(ctx context.Context, id int64) (ServiceAreaDetail, error) {
	return service.client.GetServiceArea(ctx, id)
}

func (service *OperationManagement) CreateServiceArea(ctx context.Context, input CreateServiceAreaInput) (ServiceArea, error) {
	return service.client.CreateServiceArea(ctx, input)
}

func (service *OperationManagement) UpdateServiceArea(ctx context.Context, input UpdateServiceAreaInput) (ServiceArea, error) {
	return service.client.UpdateServiceArea(ctx, input)
}

func (service *OperationManagement) ChangeServiceAreaStatus(ctx context.Context, input ChangeServiceAreaStatusInput) (ServiceArea, error) {
	return service.client.ChangeServiceAreaStatus(ctx, input)
}

func (service *OperationManagement) ReplaceServiceAreaSchedules(ctx context.Context, input ReplaceSchedulesInput) (ServiceAreaDetail, error) {
	return service.client.ReplaceServiceAreaSchedules(ctx, input)
}

func (service *OperationManagement) ListStations(ctx context.Context, input ListStationsInput) ([]Station, error) {
	return service.client.ListStations(ctx, input)
}

func (service *OperationManagement) CreateStation(ctx context.Context, input CreateStationInput) (Station, error) {
	return service.client.CreateStation(ctx, input)
}

func (service *OperationManagement) UpdateStation(ctx context.Context, input UpdateStationInput) (Station, error) {
	return service.client.UpdateStation(ctx, input)
}

func (service *OperationManagement) ChangeStationStatus(ctx context.Context, input ChangeStationStatusInput) (Station, error) {
	return service.client.ChangeStationStatus(ctx, input)
}
