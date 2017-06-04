package main

import (
	"net/http"

	"fmt"

	"os"
	"strconv"

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
	devMode, err := strconv.ParseBool(os.Getenv(config.DevMode))
	if err != nil {
		devMode = false
	}

	apiClient := halapi.HalClient{}

	logger := getLogger(devMode)

	apiClient.Log = logger
	apiClient.Router = router.HalRouter{RequestHandler: http.NewServeMux()}

	cli, err := client.NewEnvClient()
	if err != nil {
		apiClient.Log.Error(err.Error())
		return
	}
	apiClient.Engine = compose.Docker{Client: cli}
	apiClient.CreateRoutes()
	fmt.Printf("Running Hal client on port 3000")
	apiClient.Log.Info("Running Hal client on port 3000")

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

func getLogger(devMode bool) logger.SystemLogger {
	logs := new(logger.SystemLogger)
	if devMode {
		println("Running in dev mode")
		fileLogger := logger.FileLogSave{
			LogFileName: logFileName,
			LogFileDir:  os.Getenv(config.LogPath),
		}
		logs.LogSave = fileLogger
	} else {
		papertrailLog := papertrail.PapertrailLogSave{
			App: "Hal",
			URL: os.Getenv(config.LogEndpoint),
		}
		logs.LogSave = papertrailLog
	}

	return *logs
}
