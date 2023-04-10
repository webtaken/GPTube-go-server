package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/firebase_services"
	"server/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerReq models.Register
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
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

	if registerReq.Email == "" {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("please provide an email").Error(),
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

	// Adding the lead
	err := firebase_services.AddLead(registerReq.Email)
	if err != nil {
		ErrorResponse := ErrorResponseYoutube{
			ErrorResponse: fmt.Errorf("couldn't add your email").Error(),
		}
		data, err := json.Marshal(ErrorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			log.Printf("JSON marshaling failed: %s", err)
			return
		}
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
}
