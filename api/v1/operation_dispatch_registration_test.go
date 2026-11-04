package v1

import (
	"reflect"
	"testing"
)

func TestApiGroupRegistersOperationDispatchApi(t *testing.T) {
	groupType := reflect.TypeFor[ApiGroup]()
	field, ok := groupType.FieldByName("OperationDispatchApi")
	if !ok {
		t.Fatal("ApiGroup must register OperationDispatchApi")
	}

	if field.Type != reflect.TypeFor[OperationDispatchApi]() {
		t.Fatalf("OperationDispatchApi field type = %v, want %v", field.Type, reflect.TypeFor[OperationDispatchApi]())
	}
}
