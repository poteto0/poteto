package poteto

import (
	"testing"

	"github.com/poteto0/poteto/constant"
)

func TestAddAndGetParam(t *testing.T) {
	hp := NewHttpParam()

	pu := ParamUnit{"key", "value"}
	hp.AddParam(constant.PARAM_TYPE_PATH, pu)

	tests := []struct {
		name         string
		key          string
		expected_val string
		expected_ok  bool
	}{
		{"test ok case", "key", "value", true},
		{"test unexpected", "unexpected", "", false},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			value, ok := hp.GetParam(constant.PARAM_TYPE_PATH, it.key)

			if value != it.expected_val {
				t.Errorf("Don't Work")
			}

			if ok != it.expected_ok {
				t.Errorf("Unmatched")
			}
		})
	}
}

func TestJsonSerializeHttpParam(t *testing.T) {
	hp := NewHttpParam()
	hp.AddParam(constant.PARAM_TYPE_PATH, ParamUnit{key: "key", value: "value"})

	expected := `{"path":{"key":"value"},"query":{}}`
	serialized, _ := hp.JsonSerialize()
	if string(serialized) != expected {
		t.Errorf(
			"Unmatched actual(%s) -> expected(%s)",
			string(serialized),
			expected,
		)
	}
}
