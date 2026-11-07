package initialize

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// TestCsRouterPrefix 回归：客服接口必须注册为 /cs/...，不能带 /api。
// 前端 VITE_BASE_API=/api，Vite 代理会 rewrite 掉 /api 后再转发，
// 后端若注册成 /api/cs/...，前端请求会全部 404。
func TestCsRouterPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	private := engine.Group("")
	public := engine.Group("")
	initBizRouter(private, public)

	routes := make(map[string]bool)
	for _, r := range engine.Routes() {
		routes[r.Method+" "+r.Path] = true
	}

	want := []string{
		"GET /cs/session/list",
		"GET /cs/ticket/list",
		"GET /cs/faq/list",
		"GET /cs/notification/list",
	}
	for _, w := range want {
		if !routes[w] {
			t.Errorf("客服接口未注册：%s", w)
		}
	}
	if routes["GET /api/cs/session/list"] {
		t.Error("客服接口仍注册在 /api/cs 前缀下，会导致前端 404")
	}
}

// TestGvaPgDatasource 回归：系统库与业务表统一使用 db-list[gvaPg]，
// alias-name/type/db-name 必须正确，否则启动注入与业务连接都会失败。
func TestGvaPgDatasource(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("../config.dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("读取 config.dev.yaml 失败: %v", err)
	}
	list := v.Get("db-list")
	items, ok := list.([]interface{})
	if !ok {
		t.Fatal("db-list 配置缺失或格式错误")
	}
	found := false
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if m["alias-name"] == "gvaPg" {
			found = true
			if m["type"] != "pgsql" {
				t.Errorf("gvaPg.type = %v，期望 pgsql", m["type"])
			}
			if m["db-name"] != "gva" {
				t.Errorf("gvaPg.db-name = %v，期望 gva", m["db-name"])
			}
			if m["disable"] == true {
				t.Error("gvaPg 数据源被 disable")
			}
		}
	}
	if !found {
		t.Error("db-list 中缺少 alias-name: gvaPg 的数据源")
	}
}
