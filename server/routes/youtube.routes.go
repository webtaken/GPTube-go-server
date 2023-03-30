package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/YoutubeAnalyzer"
	"server/models"
)

func YoutubeHandler(w http.ResponseWriter, r *http.Request) {
	var youtubeAnalyzerReq models.YoutubeAnalyzerRequestBody

	if err := json.NewDecoder(r.Body).Decode(&youtubeAnalyzerReq); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	comments, err := YoutubeAnalyzer.GetComments(youtubeAnalyzerReq)
	if err != nil && len(comments) == 0 {
		log.Fatalf("error: %v\n", err)
	}

	data, err := json.Marshal(comments)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("Number of comments analyzed: %d\n", len(comments))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
