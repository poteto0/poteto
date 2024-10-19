package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poteto0/poteto"
	"github.com/poteto0/poteto/constant"
)

func TestCamaraWithConfigByDefault(t *testing.T) {
	config := DefaultCamaraConfig

	t.Run("allow check header", func(t *testing.T) {
		camara := CamaraWithConfig(config)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		context := poteto.NewContext(w, req)

		handler := func(ctx poteto.Context) error {
			return ctx.JSON(http.StatusOK, TestVal{Name: "test", Val: "val"})
		}

		camara_handler := camara(handler)
		camara_handler(context)
		header := w.Result().Header

		if header[constant.CONTENT_SECURITY_POLICY][0] != config.ContentSecurityPolicy {
			t.Errorf("Cannot set CSP")
		}

		if header[constant.X_FRAME_OPTION][0] != config.XFrameOption {
			t.Errorf("Cannot set XFO")
		}

		if header[constant.STRICT_TRANSPORT_SECURITY][0] != config.StrictTransportSecurity {
			t.Errorf("Cannot set STS")
		}

		if header[constant.X_DOWNLOAD_OPTION][0] != config.XDownloadOption {
			t.Errorf("Cannot set XDO")
		}

		if header[constant.X_CONTENT_TYPE_OPTION][0] != config.XContentTypeOption {
			t.Errorf("Cannot set XCT")
		}

		if header[constant.REFERRER_POLICY][0] != config.ReferrerPolicy {
			t.Errorf("Cannot set RP")
		}
	})
}
