package halapi

import (
	"encoding/json"
	"net/http"

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
		image := r.FormValue("image")
		if err := api.Engine.RunContainer(image); err != nil {
			w.Header().Set("Content-Type", "application/json")
			m := Response{Error: compose.ContainerNotFound}

			api.Log.Error(err.Error())
			resp, _ := json.Marshal(m)
			w.Write(resp)
			return
		}
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
