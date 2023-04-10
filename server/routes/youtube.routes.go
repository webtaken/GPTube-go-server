package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/YoutubeAnalyzer"
	"server/firebase_services"
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
			log.Printf("JSON marshaling failed: %s", err)
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

	videoData, err := YoutubeAnalyzer.CanProcessVideo(youtubeAnalyzerReq)
	if err != nil {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("%v", err).Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		w.WriteHeader(http.StatusBadRequest)
		if err != nil {
			log.Printf("JSON marshaling failed: %s", err)
			return
		}
		w.Write(data)
		return
	}

	// Adding lead email to temporal database
	go func() {
		firebase_services.AddLead(youtubeAnalyzerReq.Email)
	}()

	// Calling AI worker
	go func() {
		commentsResults, err := YoutubeAnalyzer.GetComments(youtubeAnalyzerReq)
		if err != nil {
			// Sending the e-mail error to the user
			subjectEmail := fmt.Sprintf(
				"GPTube analysis for YT video %q failed 😔",
				videoData.Items[0].Snippet.Title,
			)
			log.Printf("%v\n", err.Error())
			go web.SendYoutubeErrorTemplate(subjectEmail, []string{youtubeAnalyzerReq.Email})
			return
		}

		commentsResults.VideoID = youtubeAnalyzerReq.VideoID

		// Sending the e-mail to the user
		subjectEmail := fmt.Sprintf(
			"GPTube analysis for YT video %q ready 😺!",
			videoData.Items[0].Snippet.Title,
		)
		go web.SendYoutubeTemplate(
			*commentsResults, subjectEmail, []string{youtubeAnalyzerReq.Email})

		fmt.Printf("Number of comments analyzed: %d\n", commentsResults.TotalCount)
	}()

	w.WriteHeader(http.StatusOK)
}
