package halapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/Gamebuildr/Hal/pkg/auth"
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

// Response to send back to requesting service
type Response struct {
	Error          string
	ContainerCount int
}

// CreateRoutes adds specific route endpoints
func (api *HalClient) CreateRoutes() {
	api.Router.AddRoute(RunContainerRoute, api.runContainerHandler())
	api.Router.AddRoute(ContainerCountRoute, auth.JWTAuthMiddleware.Handler(api.containerCountHandler()))
}

func (api *HalClient) runContainerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		image := query.Get("image")
		if image == "" {
			api.Log.Error("No image passed in request")
			http.Error(w, "No image passed in request", http.StatusInternalServerError)
		}
		rawData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := string(rawData)
		api.Log.Info(fmt.Sprintf("Attemping to run container for %v with %v", image, data))
		if err := api.Engine.RunContainer(data, image); err != nil {
			api.Log.Error(err.Error())
			w.Header().Set("Content-Type", "application/json")
			m := Response{Error: compose.ContainerNotFound}

			api.Log.Error(err.Error())
			resp, _ := json.Marshal(m)
			w.Write(resp)
			return
		}
		api.Log.Info("Run successfully")
		w.Write([]byte(image))
	})
}

func (api *HalClient) containerCountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count, err := api.Engine.ActiveContainers()
		if err != nil {
			api.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		m := Response{ContainerCount: count}
		resp, _ := json.Marshal(m)

		w.Write(resp)
	})
}
