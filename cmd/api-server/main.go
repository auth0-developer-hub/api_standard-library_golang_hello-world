package main

import (
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/config"
	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/helpers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	clientOriginUrl := helpers.SafeGetEnv("CLIENT_ORIGIN_URL")
	port := helpers.SafeGetEnv("PORT")
	audience := helpers.SafeGetEnv("AUTH0_AUDIENCE")
	domain := helpers.SafeGetEnv("AUTH0_DOMAIN")

	config := Config{
		Port:          port,
		SecureOptions: config.SecureOptions(),
		CorsOptions:   config.CorsOptions(clientOriginUrl),
		Audience:      audience,
		Domain:        domain,
	}

	app := App{Config: config}

	app.RunServer()
}
