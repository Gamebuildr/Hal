package router

import "net/http"

// HalRouter is the main router to use with the HalClient
type HalRouter APIRouter

// AddRoute creates a new api endpoint
func (r HalRouter) AddRoute(url string, handler http.Handler) {
	r.RequestHandler.Handle(url, handler)
}
