package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/cs"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/platform"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/trip"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	PlatformApiGroup platform.ApiGroup
	TripApiGroup     trip.ApiGroup
	CsApiGroup       cs.ApiGroup
}
