package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/YoutubeAnalyzer"
	"server/models"
)

type ErrorResponseYoutube struct {
	ErrorResponse string `json:"error"`
}

func YoutubeHandler(w http.ResponseWriter, r *http.Request) {
	var youtubeAnalyzerReq models.YoutubeAnalyzerRequestBody
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&youtubeAnalyzerReq); err != nil {
		log.Printf("JSON unmarshaling failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok, err := YoutubeAnalyzer.CanProcessVideo(youtubeAnalyzerReq)
	if !ok {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("%v", err).Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	comments, err := YoutubeAnalyzer.GetComments(youtubeAnalyzerReq)
	if err != nil {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("%v", err).Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	data, err := json.Marshal(comments)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("Number of comments analyzed: %d\n", len(comments))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
