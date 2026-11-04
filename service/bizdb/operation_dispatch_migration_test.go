package bizdb

import (
	"testing"

	tripModel "github.com/flipped-aurora/gin-vue-admin/server/model/trip"
)

func TestTripBusinessModelsIncludeOperationDispatch(t *testing.T) {
	for _, model := range tripBusinessModels() {
		if _, ok := model.(*tripModel.OperationDispatch); ok {
			return
		}
	}

	t.Fatal("trip business migrations must include OperationDispatch")
}
