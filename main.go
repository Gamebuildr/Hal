package main

import (
	"net/http"

	"github.com/Gamebuildr/Hal/halapi"
	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/docker/docker/client"
)

func main() {
	apiClient := halapi.HalClient{}
	apiClient.Log = logger.SystemLogger{}
	apiClient.Router = router.HalRouter{RequestHandler: http.NewServeMux()}

	cli, err := client.NewEnvClient()
	if err != nil {
		apiClient.Log.Error(err.Error())
	}
	apiClient.Engine = compose.Docker{Client: cli}
	apiClient.CreateRoutes()
	http.ListenAndServe(":3000", apiClient.Router.RequestHandler)
}
