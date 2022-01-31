package main

import (
	"log"
	"net/http"
	"os"

	"github.com/auth0-blog/hello-golang-api/router"
    "github.com/auth0-blog/hello-golang-api/middleware"

	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
	"github.com/rs/cors"
)

func safeGetEnv(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("The environment variable '%s' doesn't exist or is not set", key)
	}
	return os.Getenv(key)
}

func main() {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error while reading the .env file %s", err)
		}
	}

	router := router.Router()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{safeGetEnv("CLIENT_ORIGIN_URL")},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge: 86400,
	})
	routerWithCORS := corsMiddleware.Handler(router)

	secureMiddleware := secure.New(secure.Options{
        STSSeconds:            	31536000,
        STSIncludeSubdomains:  	true,
        STSPreload:            	true,
        FrameDeny:             	true,
		ForceSTSHeader:			true,
        ContentTypeNosniff:    	true,
        BrowserXssFilter:      	true,
		CustomBrowserXssValue:	"0",
        ContentSecurityPolicy: 	"default-src 'self', frame-ancestors 'none'",
    })
	routerWithSecurityHeaders := secureMiddleware.Handler(routerWithCORS)

	finalHandler := middleware.HandleCacheControl(routerWithSecurityHeaders)

	server := &http.Server{
		Addr:    ":" + safeGetEnv("PORT"),
		Handler: finalHandler,
	}

	log.Printf("API server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
