package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/souravlayek/file-hosting/internal/api/handler"
)

func handleFileRoutes(r *mux.Router) {
	r.HandleFunc("/s/{id}", handler.FileDownloadHandler).Methods("GET")
	r.HandleFunc("/api/upload", handler.FileUploadHandler).Methods("POST")

}

// NewRouter returns a new router instance
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	// serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	// api routes
	handleFileRoutes(r)
	return r
}
