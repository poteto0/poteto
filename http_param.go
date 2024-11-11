package poteto

import (
	"github.com/poteto0/poteto/constant"
)

type ParamUnit struct {
	key   string
	value string
}

type httpParam struct {
	params map[string]map[string]string
}

type HttpParam interface {
	GetParam(paramType, key string) (string, bool)
	AddParam(paramType string, paramUnit ParamUnit)
}

func NewHttpParam() HttpParam {
	params := make(map[string]map[string]string, 4)

	httpParam := &httpParam{
		params: params,
	}

	httpParam.params[constant.PARAM_TYPE_PATH] = make(map[string]string)
	httpParam.params[constant.PARAM_TYPE_QUERY] = make(map[string]string)
	return httpParam
}

func (hp *httpParam) GetParam(paramType, key string) (string, bool) {
	val := hp.params[paramType][key]
	if val != "" {
		return val, true
	}
	return "", false
}

func (hp *httpParam) AddParam(paramType string, paramUnit ParamUnit) {
	hp.params[paramType][paramUnit.key] = paramUnit.value
}
