package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/cs"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/module3"
	"github.com/flipped-aurora/gin-vue-admin/server/router/platform"
	"github.com/flipped-aurora/gin-vue-admin/server/router/ride"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/trip"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Platform platform.RouterGroup
	Ride     ride.RouterGroup
	Trip     trip.RouterGroup
	Cs       cs.RouterGroup
	Module3  module3.RouterGroup

	PricingRuleRouter
	ServiceAreaRouter
	StationRouter
	OperationDispatchRouter
}
