package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInstanceTokenAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		configured    string
		header        string
		expectedCode  int
		expectedBody  string
		handlerCalled bool
	}{
		{
			name:          "disabled when not configured",
			expectedCode:  http.StatusOK,
			expectedBody:  "ok",
			handlerCalled: true,
		},
		{
			name:         "requires token when configured",
			configured:   "secret",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "rejects wrong token",
			configured:   "secret",
			header:       "wrong",
			expectedCode: http.StatusForbidden,
		},
		{
			name:          "accepts matching token",
			configured:    "secret",
			header:        "secret",
			expectedCode:  http.StatusOK,
			expectedBody:  "ok",
			handlerCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			router := gin.New()
			router.POST("/channel", InstanceTokenAuth(tt.configured), func(c *gin.Context) {
				called = true
				c.String(http.StatusOK, "ok")
			})

			request := httptest.NewRequest(http.MethodPost, "/channel", nil)
			if tt.header != "" {
				request.Header.Set(InstanceTokenHeader, tt.header)
			}
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != tt.expectedCode {
				t.Fatalf("expected status %d, got %d", tt.expectedCode, response.Code)
			}
			if tt.expectedBody != "" && response.Body.String() != tt.expectedBody {
				t.Fatalf("expected body %q, got %q", tt.expectedBody, response.Body.String())
			}
			if called != tt.handlerCalled {
				t.Fatalf("expected handlerCalled=%v, got %v", tt.handlerCalled, called)
			}
		})
	}
}
