package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/Gamebuildr/Hal/halapi"
	"github.com/Gamebuildr/Hal/pkg/channels"
	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/braintree/manners"
	"github.com/docker/docker/client"
)

const logFileName string = "hal_api_client_"

func main() {
	apiClient := halapi.HalClient{}

	// create log directory
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	logDir := path.Join(rootDir, "halapi/logs", logFileName)
	fileLogger := logger.FileLogSave{LogFileDir: logDir}
	apiClient.Log = logger.SystemLogger{LogSave: fileLogger}
	apiClient.Router = router.HalRouter{RequestHandler: http.NewServeMux()}

	cli, err := client.NewEnvClient()
	if err != nil {
		apiClient.Log.Error(err.Error())
		return
	}
	apiClient.Engine = compose.Docker{Client: cli}
	apiClient.CreateRoutes()

	apiClient.Log.Info("Running Hal Client on port 3000")
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
