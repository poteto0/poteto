package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poteto0/poteto"
)

type jwtUserClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type user struct {
	name string `json:"name"`
}

func generateUserClaims(user user, time_duration time.Duration) *jwtUserClaims {
	return &jwtUserClaims{
		user.name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time_duration)),
		},
	}
}

func TestJWSMiddleware(t *testing.T) {
	jwsConfig := NewPotetoJWSConfig(
		"user", []byte("secret"),
	)

	p := poteto.New()
	p.Register(JWSWithConfig(jwsConfig))

	p.GET("/users", func(ctx poteto.Context) error {
		return ctx.JSON(200, "hello")
	})

	u := user{
		name: "hello",
	}

	// create jws token
	claims := generateUserClaims(u, time.Hour*14*24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, _ := token.SignedString([]byte("secret"))

	tests := []struct {
		name         string
		withToken    bool
		isReplace    bool
		expectedCode int
	}{
		{"Test valid token", true, false, http.StatusOK},
		{"Test not included token", false, false, http.StatusBadRequest},
		{"Test invalid token", true, true, http.StatusUnauthorized},
	}

	for _, it := range tests {
		t.Run(it.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users", nil)

			if it.isReplace {
				tk = "invalid token"
			}

			if it.withToken {
				bearer := "Bearer " + tk
				req.Header.Set("Authorization", bearer)
			}

			p.ServeHTTP(w, req)
			if w.Code != it.expectedCode {
				fmt.Println(w.Code)
				t.Errorf("Unmatched")
			}
		})
	}
}
