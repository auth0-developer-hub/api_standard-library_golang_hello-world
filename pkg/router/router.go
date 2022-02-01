package router

import (
	"net/http"

	"github.com/auth0-blog/hello-golang-api/pkg/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/messages/public", middleware.PublicApiHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/messages/protected", middleware.ProtectedApiHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/messages/admin", middleware.AdminApiHandler).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	return router
}
