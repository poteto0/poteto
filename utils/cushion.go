package utils

import (
	"errors"
	"reflect"

	"gopkg.in/yaml.v3"
)

func YamlParse(source string, dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not pointer")
	}

	parseError := yaml.Unmarshal([]byte(source), dest)
	return parseError
}
