package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

type Metadata struct {
	Api 	string `json:"api"`
	Branch 	string `json:"branch"`
}

type ApiResponse struct {
	Metadata Metadata `json:"metadata"`
	Text string `json:"text"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

var (
	metadata			= Metadata{"api_standard-library_golang_hello-world", "starter"}
	publicMessage    	= ApiResponse{metadata, "This is a public message."}
	protectedMessage 	= ApiResponse{metadata, "This is a protected message."}
	adminMessage     	= ApiResponse{metadata, "This is an admin message."}
	notFoundMessage		= ErrorMessage{"Not Found"}
)

func safeGetEnv(key string) string {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error while reading the .env file %s", err)
		}
	}
	if os.Getenv(key) == "" {
		log.Fatalf("The environment variable '%s' doesn't exist or is not set", key)
	}
	return os.Getenv(key)
}

func publicApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, publicMessage)
}

func protectedApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, protectedMessage)
}

func adminApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, adminMessage)
}

func sendMessage(rw http.ResponseWriter, r *http.Request, data ApiResponse) {
	rw.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Print("json conversion error", err)
		return
	}
	_, err = rw.Write(bytes)
	if err != nil {
		log.Print("http response write error", err)
	}
}

func handleCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Add("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Add("Pragma", "no-cache")
		headers.Add("Expires", "0")
		next.ServeHTTP(rw, req)
		return
	})
}

func notFoundHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotFound)
	jsonResp, err := json.Marshal(notFoundMessage)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
	return
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/messages/public", publicApiHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/messages/protected", protectedApiHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/messages/admin", adminApiHandler).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

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

	finalHandler := handleCacheControl(routerWithSecurityHeaders)

	server := &http.Server{
		Addr:    ":" + safeGetEnv("PORT"),
		Handler: finalHandler,
	}

	log.Printf("API server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
