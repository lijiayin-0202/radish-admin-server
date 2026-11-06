package router

import (
	"reflect"
	"testing"
)

func TestRouterGroupRegistersOperationDispatchRouter(t *testing.T) {
	groupType := reflect.TypeFor[RouterGroup]()
	field, ok := groupType.FieldByName("OperationDispatchRouter")
	if !ok {
		t.Fatal("RouterGroup must register OperationDispatchRouter")
	}

	if field.Type != reflect.TypeFor[OperationDispatchRouter]() {
		t.Fatalf("OperationDispatchRouter field type = %v, want %v", field.Type, reflect.TypeFor[OperationDispatchRouter]())
	}
}
