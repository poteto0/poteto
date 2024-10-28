package utils

import (
	"testing"
)

type AB struct {
	A string `yaml:"a"`
	B string `yaml:"b"`
}

var data = `
a: test
b: hello
`

// tab not allowed
var not_expected = `
	a: *
	b: "hello"
`

func TestYamlParse(t *testing.T) {
	tests := []struct {
		name      string
		yaml_file any
		worked    bool
		expected  AB
	}{
		{"Test string yaml", data, true, AB{A: "test", B: "hello"}},
		{"Test []byte yaml", []byte(data), true, AB{A: "test", B: "hello"}},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			var ab AB

			switch asserted := any(it.yaml_file).(type) {
			case string:
				if err := YamlParse(asserted, &ab); err != nil {
					if it.worked {
						t.Errorf("Not expected Error")
					}
					return
				}

				if !it.worked {
					t.Errorf("Not occurred error")
					return
				}

				if it.expected.A != ab.A || it.expected.B != ab.B {
					t.Errorf("Not matched")
				}

			case []byte:
				if err := YamlParse(asserted, &ab); err != nil {
					if it.worked {
						t.Errorf("Not expected Error")
					}
					return
				}

				if !it.worked {
					t.Errorf("Not occurred error")
					return
				}

				if it.expected.A != ab.A || it.expected.B != ab.B {
					t.Errorf("Not matched")
				}
			}
		})
	}
}
