package main

import (
	"net/http"

	"fmt"

	"os"

	"github.com/Gamebuildr/Hal/halapi"
	"github.com/Gamebuildr/Hal/pkg/channels"
	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/config"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/papertrail"
	"github.com/braintree/manners"
	"github.com/docker/docker/client"
)

const logFileName string = "hal_api_client_"

func main() {
	apiClient := halapi.HalClient{}

	papertrailLogger := papertrail.PapertrailLogSave{
		App: "Hal",
		URL: os.Getenv(config.LogEndpoint),
	}
	apiClient.Log = logger.SystemLogger{LogSave: papertrailLogger}
	apiClient.Router = router.HalRouter{RequestHandler: http.NewServeMux()}

	cli, err := client.NewEnvClient()
	if err != nil {
		apiClient.Log.Error(err.Error())
		return
	}
	apiClient.Engine = compose.Docker{Client: cli}
	apiClient.CreateRoutes()
	fmt.Printf("Running Hal client on port 3000")

	server := manners.NewServer()
	server.Handler = apiClient.Router.RequestHandler
	server.Addr = ":3000"

	gracefulQuit := channels.OSSigChannel{}
	gracefulQuit.Log = apiClient.Log
	gracefulQuit.Server = server

	gracefulQuit.CreateChannel()
	gracefulQuit.StartChannelListener()

	if err := server.ListenAndServe(); err != nil {
		apiClient.Log.Error(err.Error())
		return
	}
}
