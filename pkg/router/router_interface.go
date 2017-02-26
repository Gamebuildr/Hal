package router

import "net/http"

// Router is the interface for implementing new routers
type Router interface {
	AddRoute(url string, handler http.Handler)
}

// APIRouter is the base implementation for routers
type APIRouter struct {
	RequestHandler *http.ServeMux
	Router         Router
}
