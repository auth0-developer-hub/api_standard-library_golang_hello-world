package main

import (
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := "6060"

	config := Config{
		Port:          port,
		SecureOptions: config.SecureOptions(),
		CorsOptions:   config.CorsOptions(),
	}

	app := App{Config: config}

	app.RunServer()
}
