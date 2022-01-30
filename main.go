package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/unrolled/secure"
	"github.com/rs/cors"
)

const corsAllowedDomain = "http://localhost:4040"

type Metadata struct {
	Api 	string `json:"api"`
	Branch 	string `json:"branch"`
}

type ApiResponse struct {
	Metadata Metadata `json:"metadata"`
	Text string `json:"text"`
}

var (
	metadata			= Metadata{"api_standard-library_golang_hello-world", "starter"}
	publicMessage    	= ApiResponse{metadata, "This is a public message."}
	protectedMessage 	= ApiResponse{metadata, "This is a protected message."}
	adminMessage     	= ApiResponse{metadata, "This is an admin message."}
)

func publicApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, publicMessage)
}

func protectedApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, protectedMessage)
}

func adminApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, adminMessage)
}

func sendMessage(rw http.ResponseWriter, data ApiResponse) {
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

func main() {	
	router := http.NewServeMux()
	router.Handle("/", http.NotFoundHandler())
	router.Handle("/api/messages/public", http.HandlerFunc(publicApiHandler))
	router.Handle("/api/messages/protected", http.HandlerFunc(protectedApiHandler))
	router.Handle("/api/messages/admin", http.HandlerFunc(adminApiHandler))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4040"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge: 86400,
	})
	routerWithCORS := c.Handler(router)
	secureMiddleware := secure.New(secure.Options{
        STSSeconds:            31536000,
        STSIncludeSubdomains:  	true,
        STSPreload:            	true,
        FrameDeny:             	true,
		ForceSTSHeader:			true,
        ContentTypeNosniff:    	true,
        BrowserXssFilter:      	true,
		CustomBrowserXssValue:	"0",
        ContentSecurityPolicy: 	"default-src 'self', frame-ancestors 'none'",
    })
	routerWithCacheControl := handleCacheControl(routerWithCORS)
	finalHandler := secureMiddleware.Handler(routerWithCacheControl)

	server := &http.Server{
		Addr:    ":6060",
		Handler: finalHandler,
	}

	log.Printf("API server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
