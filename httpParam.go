package poteto

import (
	"github.com/poteto0/poteto/constant"
)

type ParamUnit struct {
	key   string
	value any
}

type httpParam struct {
	params map[string]map[string]any
}

type HttpParam interface {
	GetParam(paramType, key string) any
	AddParam(paramType string, paramUnit ParamUnit)
}

func NewHttpParam() HttpParam {
	params := make(map[string]map[string]any, 4)

	httpParam := &httpParam{
		params: params,
	}

	httpParam.params[constant.PARAM_TYPE_PATH] = make(map[string]any)
	return httpParam
}

func (hp *httpParam) GetParam(paramType, key string) any {
	return hp.params[paramType][key]
}

func (hp *httpParam) AddParam(paramType string, paramUnit ParamUnit) {
	hp.params[paramType][paramUnit.key] = paramUnit.value
}
