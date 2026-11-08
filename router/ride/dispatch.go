package ride

import (
	"github.com/gin-gonic/gin"
)

type DispatchRouter struct{}

func (s *DispatchRouter) InitDispatchRouter(Router *gin.RouterGroup) {
	dispatchRouter := Router.Group("ride/dispatch")
	{
		dispatchRouter.GET("list", dispatchApi.ListDispatchRecords)
	}
}
