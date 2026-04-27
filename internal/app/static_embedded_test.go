package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/static"
)

func TestStaticFileHandlerServesEmbeddedExportedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	site := newStaticSite(static.GetWebFS())
	site.Register(router)
	router.NoRoute(site.Handle)

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		for _, target := range []string{"/favorites", "/history"} {
			request := httptest.NewRequest(method, target, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("expected %s %s status %d, got %d", method, target, http.StatusOK, response.Code)
			}
		}
	}
}
