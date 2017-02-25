package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/docker/docker/api/types"
)

// OkResponseContainerServer is a mock server to test the docker api against
// OkResponseContainerServer returns a list of containers that can be specified
func OkResponseContainerServer(content []types.Container) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		encode := encoder.Encode(content)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, encode)
	}))
}
