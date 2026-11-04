package service

import (
	"reflect"
	"testing"
)

func TestServiceGroupRegistersOperationDispatchService(t *testing.T) {
	groupType := reflect.TypeFor[ServiceGroup]()
	field, ok := groupType.FieldByName("OperationDispatchService")
	if !ok {
		t.Fatal("ServiceGroup must register OperationDispatchService")
	}

	if field.Type != reflect.TypeFor[OperationDispatchService]() {
		t.Fatalf("OperationDispatchService field type = %v, want %v", field.Type, reflect.TypeFor[OperationDispatchService]())
	}
}
