package poteto

import (
	"strings"

	"github.com/poteto0/poteto/constant"
)

type Binder interface {
	Bind(ctx Context, object any) error
}

type binder struct{}

func NewBinder() Binder {
	return &binder{}
}

func (*binder) Bind(ctx Context, object any) error {
	req := ctx.GetRequest()
	if req.ContentLength == 0 {
		return nil
	}

	base, _, _ := strings.Cut(
		req.Header.Get(constant.HEADER_CONTENT_TYPE), ";",
	)
	mediaType := strings.TrimSpace(base)

	switch mediaType {
	case constant.APPLICATION_JSON:
		if err := ctx.JsonDeserialize(object); err != nil {
			return err
		}
	}

	return nil
}
