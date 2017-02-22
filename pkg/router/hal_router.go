package router

import "net/http"

// HalRouter is the router implementation for Hal
type HalRouter struct {
	RequestHandler *http.ServeMux
}

// AddRoute creates a new api endpoint
func (r HalRouter) AddRoute(url string, endpoint endpoint) {
	r.RequestHandler.HandleFunc(url, endpoint)
}
