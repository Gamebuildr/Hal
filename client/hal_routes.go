package client

import "github.com/Gamebuildr/Hal/pkg/router"
import "net/http"
import "fmt"

// RunContainerEndpoint is the API route for running a container
const RunContainerEndpoint = "/api/container/run"

// ContainerCountEndpoint is the route for getting total number of running
// of running containers
const ContainerCountEndpoint = "/api/container/count"

// CreateRoutes for the hal client system
func CreateRoutes(halRouter router.Router) {
	halRouter.AddRoute(RunContainerEndpoint, runContainer)
	halRouter.AddRoute(ContainerCountEndpoint, containerCount)
}

func runContainer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Running Docker Container")
	return
}

func containerCount(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Get Container Count")
}
