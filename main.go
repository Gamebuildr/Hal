package main

import (
	"github.com/Gamebuildr/Hal/client"
	"github.com/Gamebuildr/Hal/pkg/router"

	"net/http"
)

func main() {
	appRouter := router.HalRouter{RequestHandler: http.NewServeMux()}
	client.CreateRoutes(appRouter)
}
