package middleware

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		// 超级管理员（authorityId=888）放行
		if waitUse.AuthorityId == 888 {
			c.Next()
			return
		}
		//获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 只读接口（数据大屏统计、客服中心列表）登录(JWTAuth)即可访问，
		// 无需为每个只读接口单独登记 casbin 策略，避免因 sys_apis / casbin_rule
		// 漏登记导致新建角色（如客服角色）「权限不足」。
		if act == "GET" && (strings.HasPrefix(obj, "/platform/dashboard/") || strings.HasPrefix(obj, "/cs/") || strings.HasPrefix(obj, "/trip/")) {
			c.Next()
			return
		}
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
