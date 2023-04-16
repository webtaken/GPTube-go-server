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

	"github.com/gorilla/mux"
)

type ErrorResponseYoutube struct {
	ErrorResponse string `json:"error"`
}

func YoutubePreAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body models.YoutubePreAnalyzerReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("%v", err).Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		if err != nil {
			log.Printf("JSON marshaling failed: %s", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}

	if body.VideoID == "" {
		errResp := models.YoutubePreAnalyzerRespBody{
			Err: fmt.Errorf("please provide a videoID and an email").Error(),
		}
		data, err := json.Marshal(errResp)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	videoData, err := YoutubeAnalyzer.CanProcessVideo(&body)
	if err != nil {
		errResp := models.YoutubePreAnalyzerRespBody{
			Err: fmt.Errorf("%v", err).Error(),
		}
		data, err := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		if err != nil {
			log.Printf("JSON marshaling failed: %s", err)
			return
		}
		w.Write(data)
		return
	}
	successResp := models.YoutubePreAnalyzerRespBody{
		VideoID:       body.VideoID,
		Snippet:       videoData.Items[0].Snippet,
		NumOfComments: int(videoData.Items[0].Statistics.CommentCount),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResp)
}

func YoutubeAnalyzerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var body models.YoutubeAnalyzerReqBody
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errResp := models.YoutubeAnalyzerRespBody{
			Err: fmt.Errorf("%v", err).Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		data, err := json.Marshal(errResp)
		if err != nil {
			log.Printf("JSON marshaling failed: %s", err)
		}
		w.Write(data)
		return
	}

	if body.Email == "" {
		// This means we haven´t received email hence is a short video so we do
		// all the logic here and send the response instantly to the client
		results, err := YoutubeAnalyzer.Analyze(body)
		if err != nil {
			// Sending the error to the user
			errResp := models.YoutubeAnalyzerRespBody{
				Err: fmt.Sprintf(
					"GPTube analysis for YT video %q failed 😔, try again later or contact us.",
					body.VideoTitle,
				),
			}
			w.WriteHeader(http.StatusInternalServerError)
			data, err := json.Marshal(errResp)
			if err != nil {
				log.Printf("JSON marshaling failed: %s", err)
			}
			w.Write(data)
		} else {
			// sending the results to the user
			successResp := models.YoutubeAnalyzerRespBody{
				VideoID:      vars["videoID"],
				BertAnalysis: results,
			}
		}
		return
	}

	// This means we have received email hence this video is large so we do all
	// the logic in the server and send the result back to the email of the user

	// Adding lead email to temporal database
	go func() {
		firebase_services.AddLead(body.Email)
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
