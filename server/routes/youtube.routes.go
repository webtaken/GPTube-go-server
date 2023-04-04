package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/YoutubeAnalyzer"
	"server/models"
	web "server/web/youtube"
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

	go func() {
		commentsResults, err := YoutubeAnalyzer.GetComments(youtubeAnalyzerReq)
		commentsResults.VideoID = youtubeAnalyzerReq.VideoID

		if err != nil {
			fmt.Println("$$")
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

		// Sending the e-mail to the user
		go web.SendTemplate(*commentsResults)

		fmt.Printf("Number of comments analyzed: %d\n", commentsResults.TotalCount)
	}()

	w.WriteHeader(http.StatusOK)
}
