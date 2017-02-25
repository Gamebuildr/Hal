package halapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Gamebuildr/Hal/pkg/compose"
	"github.com/Gamebuildr/Hal/pkg/router"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
)

// HalClient is the implementation of the Hal api
type HalClient struct {
	Engine compose.Engine
	Router router.HalRouter
	Log    logger.Log
}

// CreateRoutes adds specific route endpoints
func (api *HalClient) CreateRoutes() {
	api.Router.AddRoute(RunContainerRoute, api.runContainer)
	api.Router.AddRoute(ContainerCountRoute, api.containerCount)
}

func (api *HalClient) runContainer(w http.ResponseWriter, r *http.Request) {
	containerData := compose.ContainerData{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.Unmarshal(data, &containerData); err != nil {
		api.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := api.Engine.RunContainer(containerData.Image); err != nil {
		api.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *HalClient) containerCount(w http.ResponseWriter, r *http.Request) {
	count, err := api.Engine.ActiveContainers()
	if err != nil {
		api.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, count)
}
