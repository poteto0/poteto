package engine

import (
	"errors"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/poteto0/poteto/cmd/template"
)

func TestCreateAndWrite(t *testing.T) {
	defer monkey.UnpatchAll()

	tests := []struct {
		name          string
		mockCreate    func(name string) (*os.File, error)
		mockWrite     func(f *os.File, b []byte) (n int, err error)
		expectedError bool
	}{
		{
			"test cannot create file throw error",
			func(name string) (*os.File, error) {
				return &os.File{}, errors.New("error")
			},
			func(f *os.File, b []byte) (n int, err error) {
				return 1, nil
			},
			true,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			monkey.Patch(os.Create, it.mockCreate)
			monkey.Patch((*os.File).Write, it.mockWrite)

			result := createAndWrite("test", "template")
			if it.expectedError && result == nil {
				t.Error("Unmatched")
			}
			if !it.expectedError && result != nil {
				t.Error("Unmatched")
			}
		})
	}
}

func TestChoiceTemplateFile(t *testing.T) {
	tests := []struct {
		name     string
		param    EngineNewParam
		expected string
	}{
		{
			"Test choice fast api template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      true,
				IsJSONRPC:   false,
			},
			template.FastTemplate,
		},
		{
			"Test choice fast jsonrpc template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      true,
				IsJSONRPC:   true,
			},
			template.JSONRPCFastTemplate,
		},
		{
			"Test choice normal jsonrpc template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      false,
				IsJSONRPC:   true,
			},
			template.JSONRPCTemplate,
		},
		{
			"Test choice normal api template",
			EngineNewParam{
				ProjectName: "test",
				IsFast:      false,
				IsJSONRPC:   false,
			},
			template.DefaultTemplate,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := choiceTemplateFile(it.param)
			if result != it.expected {
				t.Error("Unmatched")
			}
		})
	}
}
