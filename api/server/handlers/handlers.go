package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiHandler struct{}

func (ah *apiHandler) errorResponse(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
	return
}

func (ah *apiHandler) jsonResponse(w http.ResponseWriter, object interface{}) {
	ah.jsonResponseWithStatus(w, object, http.StatusOK)
}

func (ah *apiHandler) jsonResponseWithStatus(w http.ResponseWriter, object interface{}, status int) {
	data, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		ah.errorResponse(w, ErrServerError, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
