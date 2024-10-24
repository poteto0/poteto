package utils

import (
	"errors"
	"reflect"

	"github.com/goccy/go-yaml"
)

func YamlParse[T string | []byte](source T, dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not pointer")
	}

	switch source_asserted := any(source).(type) {
	case string:
		parseError := yaml.Unmarshal([]byte(source_asserted), dest)
		return parseError
	case []byte:
		parseError := yaml.Unmarshal(source_asserted, dest)
		return parseError
	default:
		return errors.New("unexpected input")
	}
}
