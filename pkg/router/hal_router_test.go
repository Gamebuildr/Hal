package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const mockURL string = "/mock/url"
const okResponse = 200

func MockResponse(w http.ResponseWriter, r *http.Request) {
	return
}

func TestHalRouteCanAddNewRoutes(t *testing.T) {
	mockrouter := HalRouter{RequestHandler: http.NewServeMux()}
	mockrouter.AddRoute(mockURL, MockResponse)

	r, err := http.NewRequest("POST", mockURL, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockrouter.RequestHandler.ServeHTTP(w, r)

	if w.Code != okResponse {
		t.Errorf("Expected %v, but got %v", okResponse, w.Code)
	}
}
