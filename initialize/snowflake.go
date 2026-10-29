package initialize

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// InitSnowflake 初始化全局雪花ID生成器
// workerId 固定为 0（本地开发/单机部署）；多机部署时按机器编号分配 0~1023
func InitSnowflake() {
	var err error
	global.Snow, err = global.NewSnowFlake(0)
	if err != nil {
		panic(fmt.Sprintf("雪花ID生成器初始化失败: %v", err))
	}
	zap.L().Info("雪花ID生成器初始化成功")
}
