package poteto

import (
	stdContext "context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/ybbus/jsonrpc/v3"
)

type (
	TestCalculator struct{}
	AdditionArgs   struct {
		Add, Added int
	}
)

func (tc *TestCalculator) Add(r *http.Request, args *AdditionArgs) int {
	return args.Add + args.Added
}

func TestPotetoJSONRPCAdapterCall(t *testing.T) {
	p := New()

	rpc := TestCalculator{}
	p.POST("/add", func(ctx Context) error {
		return PotetoJsonRPCAdapter[TestCalculator, AdditionArgs](ctx, &rpc)
	})

	errChan := make(chan error)
	go func() {
		errChan <- p.Run("6000")
	}()

	// client
	added := 10
	add := 10
	rpcClient := jsonrpc.NewClient("http://localhost:6000/add")
	//result := &AdditionResult{}
	result, err := rpcClient.Call(stdContext.Background(), "TestCalculator.Add", &AdditionArgs{Added: added, Add: add})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	num, _ := result.Result.(json.Number).Int64()
	if int(num) != 20 {
		t.Errorf("Unmatched actual(%v) -> expected(%v)", result.Result, 20)
	}

	select {
	case <-time.After(500 * time.Millisecond):
		if err := p.Stop(stdContext.Background()); err != nil {
			t.Errorf("Unmatched")
		}
	case <-errChan:
		t.Errorf("Unexpected error occur")
	}
}

func TestPotetoJSONRPCAdapter(t *testing.T) {
	w := httptest.NewRecorder()

	tests := []struct {
		name     string
		req      *http.Request
		expected float64
	}{
		{
			"Test not POST req",
			httptest.NewRequest("GET", "/test", nil),
			-32700,
		},
		{
			"Test not POST body nil",
			httptest.NewRequest("POST", "/test", nil),
			-32700,
		},
		{
			"Test not POST body not version right",
			httptest.NewRequest("POST", "/test", strings.NewReader("1")),
			-32700,
		},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			ctx := NewContext(w, it.req).(*context)

			rpc := TestCalculator{}
			PotetoJsonRPCAdapter[TestCalculator, int](ctx, &rpc)

			var data map[string]any
			json.Unmarshal(w.Body.Bytes(), &data)
			err := data["error"].(map[string]any)
			if err["code"] != it.expected {
				t.Errorf(
					"Unmatched actual(%v) -> expected(%v)",
					err["code"],
					it.expected,
				)
			}
		})
	}
}
