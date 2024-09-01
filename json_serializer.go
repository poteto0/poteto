package poteto

import "encoding/json"

type JsonSerializer interface {
	Serialize(c Context, value any) error
}

type jsonSerializer struct{}

func (j jsonSerializer) Serialize(ctx Context, value any) error {
	encoder := json.NewEncoder(ctx.GetResponse())
	return encoder.Encode(value)
}
