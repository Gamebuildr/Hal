package compose

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func okResponseTestServer(content []types.Container) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		encode := encoder.Encode(content)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, encode)
	}))
}

func TestDockerCanRunSpecifiedContainer(t *testing.T) {
	mockImage := "mock/image:latest"
	mockID := "123456789"
	mockContainers := []types.Container{
		types.Container{ID: mockID, Image: mockImage},
		types.Container{ID: "987654321", Image: "test/mock/images:latest"},
	}

	testServer := okResponseTestServer(mockContainers)
	defer testServer.Close()

	cli, err := client.NewClient(testServer.URL, "1.13", nil, map[string]string{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	engine := Docker{client: cli}
	container := Container{Engine: engine}

	if err := container.Engine.RunContainer(mockImage); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDockerRunContainerReturnsErrorWhenContainerNotFound(t *testing.T) {
	mockImage := "mock/image:latest"
	mockID := "123456789"
	mockContainers := []types.Container{
		types.Container{ID: mockID, Image: mockImage},
		types.Container{ID: "987654321", Image: "test/mock/images:latest"},
	}

	testServer := okResponseTestServer(mockContainers)
	defer testServer.Close()

	cli, err := client.NewClient(testServer.URL, "1.13", nil, map[string]string{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	engine := Docker{client: cli}
	container := Container{Engine: engine}

	runErr := container.Engine.RunContainer("no_container")
	if runErr == nil {
		t.Errorf("Expected error container not found")
	}
}

func TestDockerReturnsRunningContainerCount(t *testing.T) {
	mockContainers := []types.Container{
		types.Container{ID: "test_one"},
		types.Container{ID: "test_two"},
	}

	testServer := okResponseTestServer(mockContainers)
	defer testServer.Close()

	cli, err := client.NewClient(testServer.URL, "1.13", nil, map[string]string{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	engine := Docker{client: cli}
	container := Container{Engine: engine}

	count, err := container.Engine.ActiveContainers()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if count == 0 {
		t.Errorf("Expected 2 continers, got %v", count)
	}
	if count != 2 {
		t.Errorf("Expected 2 containers, got %v", count)
	}
}

func TestDockerCanGetContainerIDByImageName(t *testing.T) {
	mockImage := "mock/image:latest"
	mockID := "123456789"
	mockContainers := []types.Container{
		types.Container{ID: mockID, Image: mockImage},
		types.Container{ID: "987654321", Image: "test/mock/images:latest"},
	}
	testServer := okResponseTestServer(mockContainers)
	defer testServer.Close()

	cli, err := client.NewClient(testServer.URL, "1.13", nil, map[string]string{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	engine := Docker{client: cli}
	container := Container{Engine: engine}

	id, err := container.Engine.getContainerID(mockImage)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if id == "" {
		t.Errorf("Expected container id to not be empty")
	}
	if id != mockID {
		t.Errorf("Expected container id %v, got %v", mockID, id)
	}
}
