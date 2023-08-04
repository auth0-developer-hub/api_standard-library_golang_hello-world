package main

import (
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/router"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

type Config struct {
	Port          string
	SecureOptions secure.Options
	CorsOptions   cors.Options
}

type App struct {
	Config Config
}

func (app *App) RunServer() {
	router := router.Router()
	corsMiddleware := cors.New(app.Config.CorsOptions)
	routerWithCORS := corsMiddleware.Handler(router)

	secureMiddleware := secure.New(app.Config.SecureOptions)
	finalHandler := secureMiddleware.Handler(routerWithCORS)

	server := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: finalHandler,
	}

	log.Printf("API server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
