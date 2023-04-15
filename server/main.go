package main

import (
	"fmt"
	"log"
	"net/http"

	envManager "server/env_manager"
	"server/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{envManager.GoDotEnvVariable("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r.HandleFunc("/", routes.HomeHandler).Methods("GET")
	r.HandleFunc("/YT", routes.YoutubePreAnalysisHandler).Methods("POST")
	// r.HandleFunc("/YT/{videoID}", routes.YoutubeAnalyzerHandler).Methods("POST")
	r.HandleFunc("/register", routes.RegisterHandler).Methods("POST")
	port := fmt.Sprintf(":%s", envManager.GoDotEnvVariable("PORT"))
	log.Fatal(http.ListenAndServe(
		port,
		handlers.CORS(originsOk, headersOk, methodsOk)(r),
	))
}
