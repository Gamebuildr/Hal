package client

import (
	"net/http"

	"testing"

	"net/http/httptest"

	"github.com/Gamebuildr/Hal/pkg/router"
)

const okResponse = 200

func TestRouteCreationAddsRoutesCorrectly(t *testing.T) {
	mockrouter := router.HalRouter{RequestHandler: http.NewServeMux()}
	CreateRoutes(mockrouter)

	r, err := http.NewRequest("POST", RunContainerEndpoint, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockrouter.RequestHandler.ServeHTTP(w, r)

	if w.Code != okResponse {
		t.Errorf("Expected %v, but got %v", okResponse, w.Code)
	}
}
