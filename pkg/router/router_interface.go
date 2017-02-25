package router

import "net/http"

type endpoint func(http.ResponseWriter, *http.Request)

// Router is the interface for implementing new routers
type Router interface {
	AddRoute(url string, endpoint endpoint)
}

// APIRouter is the base implementation for routers
type APIRouter struct {
	RequestHandler *http.ServeMux
	Router         Router
}
