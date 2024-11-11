package poteto

import (
	"net"
	"net/http/httptest"
	"testing"

	"github.com/poteto0/poteto/constant"
)

func TestCanTrust(t *testing.T) {
	iph := &ipHandler{}

	// test setter too
	iph.SetIsTrustPrivateIP(true)
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/24")
	// test register too
	iph.RegisterTrustIPRange(ipnet)

	tests := []struct {
		name     string
		ip       net.IP
		expected bool
	}{
		{"test trusted IP", net.ParseIP("10.0.0.1"), true},
		{"test untrusted IP", net.ParseIP("111.0.0.0"), false},
		{"test private IP", net.ParseIP("192.168.0.1"), true},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			result := iph.CanTrust(it.ip)
			if result != it.expected {
				t.Errorf("Not matched")
			}
		})
	}
}

func TestGetIPFromXFFHeader(t *testing.T) {
	iph := &ipHandler{}
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/24")
	iph.RegisterTrustIPRange(ipnet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(constant.HEADER_X_FORWARDED_FOR, "11.0.0.1, 12.0.0.1, 10.0.0.2, 10.0.0.1")
	ctx := NewContext(w, req).(*context)

	ipString, _ := iph.GetIPFromXFFHeader(ctx)
	if ipString != "12.0.0.1" {
		t.Errorf("Not matched")
	}
}
