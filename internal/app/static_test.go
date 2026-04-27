package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

func TestStaticFileHandlerServesExportedRoutes(t *testing.T) {
	router := newStaticTestRouter()

	response := performStaticRequest(router, "/favorites")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if response.Body.String() != "favorites page" {
		t.Fatalf("expected favorites page, got %q", response.Body.String())
	}
}

func TestRegisterExportedPageRoutesServesExtensionlessPages(t *testing.T) {
	router := gin.New()
	webFS := fstest.MapFS{
		"index.html":     {Data: []byte("index page")},
		"favorites.html": {Data: []byte("favorites page")},
		"history.html":   {Data: []byte("history page")},
	}

	registerExportedPageRoutes(router, webFS)
	router.NoRoute(StaticFileHandler(webFS))

	for _, target := range []string{"/favorites", "/history"} {
		response := performStaticRequest(router, target)

		if response.Code != http.StatusOK {
			t.Fatalf("expected GET %s status %d, got %d", target, http.StatusOK, response.Code)
		}
	}
}

func TestStaticFileHandlerLeavesAPIRoutesAlone(t *testing.T) {
	router := newStaticTestRouter()

	response := performStaticRequest(router, "/api/clipboard/favorites")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
	if response.Body.String() == "index page" {
		t.Fatalf("expected API 404 to remain untouched, got index page")
	}
}

func TestStaticFileHandlerRejectsParentPathSegments(t *testing.T) {
	router := newStaticTestRouter()

	response := performStaticRequest(router, "/../secret.txt")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
	if response.Body.String() == "secret" {
		t.Fatalf("expected parent path segment to be rejected")
	}
}

func newStaticTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.NoRoute(StaticFileHandler(fstest.MapFS{
		"index.html":     {Data: []byte("index page")},
		"favorites.html": {Data: []byte("favorites page")},
		"secret.txt":     {Data: []byte("secret")},
	}))

	return router
}

func performStaticRequest(router http.Handler, target string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, target, nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	return response
}
