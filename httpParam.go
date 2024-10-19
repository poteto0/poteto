package poteto

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

func (hp *httpParam) GetParam(paramType, key string) any {
	return hp.params[paramType][key]
}

func (hp *httpParam) AddParam(paramType string, paramUnit ParamUnit) {
	hp.params[paramType][paramUnit.key] = paramUnit.value
}
