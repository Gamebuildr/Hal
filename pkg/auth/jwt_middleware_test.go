package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gamebuildr/Hal/pkg/testutils"
)

const mockURL string = "/mock/auth"

func mockAuthRoute() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("MockResponse"))
	})
}

func TestRouteIsValidatedBySpecifiedToken(t *testing.T) {
	mockRouter := http.NewServeMux()
	mockRouter.Handle(mockURL, JWTAuthMiddleware.Handler(mockAuthRoute()))

	r, err := http.NewRequest("POST", mockURL, nil)
	testutils.AuthenticateRoute(r)
	if err != nil {
		t.Fatalf(err.Error())
	}
	w := httptest.NewRecorder()
	mockRouter.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}
}

func TestRouteIsDeniedAccessIfTokenIncorrect(t *testing.T) {
	token, err := testutils.GetMockToken("wrongkey")
	if err != nil {
		t.Fatalf(err.Error())
	}

	mockRouter := http.NewServeMux()
	mockRouter.Handle(mockURL, JWTAuthMiddleware.Handler(mockAuthRoute()))
	bearer := "Bearer " + token

	r, err := http.NewRequest("POST", mockURL, nil)
	r.Header.Add("Authorization", bearer)
	if err != nil {
		t.Fatalf(err.Error())
	}
	w := httptest.NewRecorder()
	mockRouter.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected %v, but got %v", http.StatusUnauthorized, w.Code)
	}
}
