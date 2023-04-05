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

	if youtubeAnalyzerReq.VideoID == "" || youtubeAnalyzerReq.Email == "" {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("please provide a videoID and an email").Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
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
			fmt.Printf("error here: %v\n", err.Error())
			// Sending the e-mail error to the user
			subjectEmail := fmt.Sprintf(
				"GPTube Youtube video analysis for %s failed 😔", youtubeAnalyzerReq.VideoID)
			go web.SendYoutubeErrorTemplate(subjectEmail, []string{youtubeAnalyzerReq.Email})
			return
		}

		// Sending the e-mail to the user
		subjectEmail := fmt.Sprintf(
			"GPTube Youtube video analysis for %s ready 😺!", youtubeAnalyzerReq.VideoID)
		go web.SendYoutubeTemplate(
			*commentsResults, subjectEmail, []string{youtubeAnalyzerReq.Email})

		fmt.Printf("Number of comments analyzed: %d\n", commentsResults.TotalCount)
	}()

	w.WriteHeader(http.StatusOK)
}
