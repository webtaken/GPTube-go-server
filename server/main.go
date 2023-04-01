package main

import (
	"fmt"
	"log"
	"net/http"

	"server/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var EnvVars map[string]string

func init() {
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	EnvVars = env
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", routes.HomeHandler).Methods("GET")
	r.HandleFunc("/YT", routes.YoutubeHandler).Methods("POST")
	port := fmt.Sprintf(":%s", EnvVars["PORT"])
	http.ListenAndServe(port, r)
}
