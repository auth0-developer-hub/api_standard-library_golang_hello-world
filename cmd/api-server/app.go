package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/router"

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

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("API server closed: err: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("got shutdown signal. shutting down server...")

	localCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(localCtx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("server shutdown complete")
}
