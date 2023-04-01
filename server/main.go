package main

import (
	"fmt"
	"net/http"

	envManager "server/env_manager"
	"server/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", routes.HomeHandler).Methods("GET")
	r.HandleFunc("/YT", routes.YoutubeHandler).Methods("POST")
	port := fmt.Sprintf(":%s", envManager.GoDotEnvVariable("PORT"))
	http.ListenAndServe(port, r)
}
