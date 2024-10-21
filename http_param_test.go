package poteto

import (
	"testing"

	"github.com/poteto0/poteto/constant"
)

func TestAddAndGetParam(t *testing.T) {
	hp := NewHttpParam()

	pu := ParamUnit{"key", "value"}
	hp.AddParam(constant.PARAM_TYPE_PATH, pu)

	value := hp.GetParam(constant.PARAM_TYPE_PATH, "key")
	if value != "value" {
		t.Errorf("Don't Work")
	}
}
