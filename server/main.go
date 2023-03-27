package main

import (
	"net/http"

	"server/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", routes.HomeHandler).Methods("GET")
	r.HandleFunc("/YT", routes.YoutubeHandler).Methods("POST")
	http.ListenAndServe(":8000", r)
}
