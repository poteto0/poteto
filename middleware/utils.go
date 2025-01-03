package middleware

import (
	"regexp"
	"strings"

	"github.com/poteto-go/poteto/constant"
)

// EX: https://example.com:* => ^https://example\.com:.*$
func wrapRegExp(target string) string {
	pattern := regexp.QuoteMeta(target) // .をescapeする
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\?", ".")
	pattern = "^" + pattern + "$"
	return pattern
}

// just sub domain
// only wild card
func matchSubdomain(domain, pattern string) bool {
	if !matchScheme(domain, pattern) {
		return false
	}

	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}

	// more fast on opp
	domAuth := domain[didx+3:] // after [://]

	// avoid too long
	if len(domAuth) > constant.MAX_DOMAIN_LENGTH {
		return false
	}
	patAuth := pattern[pidx+3:]

	// Opposite by .
	domComp := strings.Split(domAuth, ".")
	domComp = reverseStringArray(domComp)
	// do pattern
	patComp := strings.Split(patAuth, ".")
	patComp = reverseStringArray(patComp)

	for i, dom := range domComp {
		if len(patComp) <= i {
			return false
		}

		pat := patComp[i]
		if pat == "*" {
			return true
		}

		if pat != dom {
			return false
		}
	}
	return false
}

// http vs https
func matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

func reverseStringArray(targets []string) []string {
	n := len(targets)
	for i := n/2 - 1; i >= 0; i-- {
		oppidx := n - i - 1
		targets[i], targets[oppidx] = targets[oppidx], targets[i]
	}

	return targets
}

func matchMethod(method string, allowMethods []string) bool {
	for _, m := range allowMethods {
		if m == method {
			return true
		}
	}

	return false
}
