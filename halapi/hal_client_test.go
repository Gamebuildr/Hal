package halapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/Hal/pkg/testutils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const mockImage = "mock/image:latest"
const containerID = "123456789"

func createMockClient(ts *httptest.Server, t *testing.T) (HalClient, error) {
	mockClient := HalClient{}

	cli, err := client.NewClient(ts.URL, "1.13", nil, map[string]string{})
	if err != nil {
		return mockClient, err
	}
	mockClient.Engine = compose.Docker{Client: cli}
	mockClient.Router = router.HalRouter{RequestHandler: http.NewServeMux()}
	mockClient.Log = testutils.MockLog{Test: t}

	return mockClient, nil
}

func TestHalClientRunContainerFindsCorrectImage(t *testing.T) {
	mockContainers := []types.Container{
		types.Container{ID: containerID, Image: mockImage},
	}
	testServer := testutils.OkResponseContainerServer(mockContainers)
	defer testServer.Close()

	mockClient, err := createMockClient(testServer, t)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockClient.CreateRoutes()

	message := []byte(`{"message": "message"}`)
	r, err := http.NewRequest("POST", RunContainerRoute+"?image=mock/image:latest", bytes.NewBuffer(message))
	testutils.AuthenticateRoute(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockClient.Router.RequestHandler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}
}

func TestHalClientRunContainerReturnsErrorWhenImageNotFound(t *testing.T) {
	mockContainers := []types.Container{
		types.Container{ID: containerID, Image: mockImage},
	}
	testServer := testutils.OkResponseContainerServer(mockContainers)
	defer testServer.Close()

	mockClient, err := createMockClient(testServer, t)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockClient.CreateRoutes()

	message := []byte(`{"message": "message"}`)
	r, err := http.NewRequest("POST", RunContainerRoute+"?image=different/image:latest", bytes.NewBuffer(message))
	testutils.AuthenticateRoute(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockClient.Router.RequestHandler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}

	jsonResp := Response{}
	resp, err := ioutil.ReadAll(w.Body)
	if err := json.Unmarshal(resp, &jsonResp); err != nil {
		t.Fatalf(err.Error())
	}

	if jsonResp.Error != compose.ContainerNotFound {
		t.Errorf("Expected %v, but got %v", compose.ContainerNotFound, jsonResp.Error)
	}
}

func TestHalClientContainerActiveCountRouteAccessible(t *testing.T) {
	testServer := testutils.OkResponseContainerServer([]types.Container{})
	defer testServer.Close()

	mockClient, err := createMockClient(testServer, t)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockClient.CreateRoutes()

	r, err := http.NewRequest("GET", ContainerCountRoute, nil)
	testutils.AuthenticateRoute(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockClient.Router.RequestHandler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}
}

func TestHalClientContainerActiveCountReturnsCorrectNumber(t *testing.T) {
	mockContainers := []types.Container{
		types.Container{ID: containerID, Image: mockImage},
		types.Container{ID: "987615321", Image: "alternative/mock:latest"},
	}
	testServer := testutils.OkResponseContainerServer(mockContainers)
	defer testServer.Close()

	mockClient, err := createMockClient(testServer, t)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockClient.CreateRoutes()

	r, err := http.NewRequest("GET", ContainerCountRoute, nil)
	testutils.AuthenticateRoute(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	w := httptest.NewRecorder()
	mockClient.Router.RequestHandler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %v, but got %v", http.StatusOK, w.Code)
	}

	jsonResp := Response{}
	resp, err := ioutil.ReadAll(w.Body)
	if err := json.Unmarshal(resp, &jsonResp); err != nil {
		t.Fatalf(err.Error())
	}
	if jsonResp.ContainerCount != 2 {
		t.Errorf("Expected %v, but got %v", 2, jsonResp.ContainerCount)
	}
}
