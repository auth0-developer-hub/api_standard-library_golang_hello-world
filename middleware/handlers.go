package middleware

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/auth0-blog/hello-golang-api/models"
)

var (
	metadata			= models.Metadata{"api_standard-library_golang_hello-world", "starter"}
	publicMessage    	= models.ApiResponse{metadata, "This is a public message."}
	protectedMessage 	= models.ApiResponse{metadata, "This is a protected message."}
	adminMessage     	= models.ApiResponse{metadata, "This is an admin message."}
	notFoundMessage		= models.ErrorMessage{"Not Found"}
)

func PublicApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, publicMessage)
}

func ProtectedApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, protectedMessage)
}

func AdminApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, adminMessage)
}

func sendMessage(rw http.ResponseWriter, r *http.Request, data models.ApiResponse) {
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

func HandleCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Add("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Add("Pragma", "no-cache")
		headers.Add("Expires", "0")
		next.ServeHTTP(rw, req)
		return
	})
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNotFound)
	jsonResp, err := json.Marshal(notFoundMessage)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
	return
}
