package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/Gamebuildr/Hal/halapi"
	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
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
	}
	apiClient.Engine = compose.Docker{Client: cli}
	apiClient.CreateRoutes()

	http.ListenAndServe(":3000", apiClient.Router.RequestHandler)
}
