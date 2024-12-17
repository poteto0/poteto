package poteto

import (
	stdContext "context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

type (
	TestCalculator struct{}
	AdditionArgs   struct {
		Add, Added int
	}
)

func (tc *TestCalculator) Add(r *http.Request, args *AdditionArgs) int {
	fmt.Println("おはようございます")
	fmt.Println(r)
	fmt.Println(args)
	return args.Add + args.Added
}

func TestJSONRPCAdapter(t *testing.T) {
	p := New()

	rpc := TestCalculator{}
	call := reflect.ValueOf(&rpc).MethodByName("Add")
	fmt.Println(call)
	p.POST("/add", func(ctx Context) error {
		return PotetoJsonRPCAdapter[TestCalculator, AdditionArgs](ctx, &rpc)
	})

	errChan := make(chan error)
	go func() {
		errChan <- p.Run("8080")
	}()

	// client
	added := 10
	add := 10
	rpcClient := jsonrpc.NewClient("http://localhost:8080/add")
	//result := &AdditionResult{}
	result, err := rpcClient.Call(stdContext.Background(), "TestCalculator.Add", &AdditionArgs{Added: added, Add: add})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	t.Errorf("%v", result)

	select {
	case <-time.After(500 * time.Millisecond):
		if err := p.Stop(stdContext.Background()); err != nil {
			t.Errorf("Unmatched")
		}
	case <-errChan:
		t.Errorf("Unexpected error occur")
	}
}
