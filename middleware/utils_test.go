package middleware

import (
	"fmt"
	"testing"
)

func TestWrapRegExp(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected string
	}{
		{"test * url", "https://example.com:*", `^https://example\.com:.*$`},
		{"test ? url", "https://example.com:300?", `^https://example\.com:300.$`},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := wrapRegExp(it.target)
			if result != it.expected {
				t.Errorf("Not matched")
			}
		})
	}
}

func TestMatchScheme(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		pattern  string
		expected bool
	}{
		{"If : is not existed return false(both)", "example.com", "example.com", false},
		{"If : is not existed return false(pattern)", "example.com", "http://example.com", false},
		{"If : is not existed return false(domain)", "http://example.com", "example.com", false},
		{"matched", "http://example1.com", "http://example2.com", true},
		{"not matched", "http://example1.com", "https://example2.com", false},
	}

	for _, it := range tests {
		t.Run((it.name), func(t *testing.T) {
			result := matchScheme(it.domain, it.pattern)

			if result != it.expected {
				t.Errorf("Not matched")
				t.Errorf(fmt.Sprintf("expected: %t", it.expected))
				t.Errorf(fmt.Sprintf("actual: %t", result))
			}
		})
	}
}
