package router

type HalRouter APIRouter

// AddRoute creates a new api endpoint
func (r HalRouter) AddRoute(url string, endpoint endpoint) {
	r.RequestHandler.HandleFunc(url, endpoint)
}
