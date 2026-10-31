package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/cs"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/platform"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/trip"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Platform platform.RouterGroup
	Trip     trip.RouterGroup
	Cs       cs.RouterGroup
}
