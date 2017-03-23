package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// OkResponseContainerServer is a mock server to test the docker api against
func OkResponseContainerServer(content []types.Container) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		encode := encoder.Encode(content)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, encode)
	}))
}

// CreateContainerServer returns allows container creation and container running to be mocked
func CreateContainerServer(content []container.ContainerCreateCreatedBody) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(content) == 0 {
			return
		}
		encoder := json.NewEncoder(w)
		encode := encoder.Encode(content[0])
		fmt.Fprintln(w, encode)
	}))
}
