package bizdb

import (
	"reflect"
	"testing"
)

// TestTripBusinessModelsCoverCoreTables 回归：统一到 db-list[gvaPg] 后，
// 核心业务表必须纳入 AutoMigrate，否则 gva 库中缺表会导致后台页面
// 报 relation does not exist、列表无数据，以及 passengers.deleted_at 列缺失。
func TestTripBusinessModelsCoverCoreTables(t *testing.T) {
	want := []string{
		"trip_orders",
		"passengers",
		"service_areas",
		"vehicles",
		"trip_status_logs",
		"safety_incidents",
	}

	got := make(map[string]bool)
	for _, model := range tripBusinessModels() {
		v := reflect.ValueOf(model)
		method := v.MethodByName("TableName")
		if !method.IsValid() {
			continue
		}
		got[method.Call(nil)[0].String()] = true
	}

	for _, table := range want {
		if !got[table] {
			t.Errorf("AutoMigrate 清单缺少核心业务表 %q", table)
		}
	}
}
