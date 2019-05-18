package server

import (
	"log"
	"net/http"
)

type emailsAPIHandler struct{}

func newEmailsAPIHandler() *emailsAPIHandler {
	return &emailsAPIHandler{}
}

func (ah *emailsAPIHandler) createEmail(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := []byte("{\"test\": \"OK\"}")

	w.WriteHeader(http.StatusOK)

	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}

	return
}
