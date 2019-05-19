package handlers

import (
	"errors"
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/parsers/emails"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/validators"
	"log"
	"net/http"
)

type EmailsAPIHandler struct {
	*apiHandler
	repository repositories.EmailsRepository
}

func NewEmailsAPIHandler(repository repositories.EmailsRepository) *EmailsAPIHandler {
	return &EmailsAPIHandler{apiHandler: &apiHandler{}, repository: repository}
}

func (eh *EmailsAPIHandler) CreateEmail(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	parser, err := emails.NewPayloadParser(contentType)
	if err != nil {
		eh.errorResponse(w, err, http.StatusUnsupportedMediaType)
		return
	}

	if err = parser.Parse(req.Body); err != nil {
		eh.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	changeSet := parser.Data()
	validator := validators.NewEmailValidator(changeSet)
	validator.Validate()

	if !validator.Valid() {
		eh.jsonResponseWithStatus(w, validator.Errors(), http.StatusUnprocessableEntity)
		return
	}

	// TODO: extract creation of an email to a service
	email := &models.Email{Status: "created"}
	changeSet.ApplyChanges(email)

	email, err = eh.repository.Create(userID(req), email)
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, email)

	return
}

func (eh *EmailsAPIHandler) ListEmails(w http.ResponseWriter, req *http.Request) {
	allEmails, err := eh.repository.All(userID(req))
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, allEmails)

	return
}

func (eh *EmailsAPIHandler) SendEmail(w http.ResponseWriter, req *http.Request) {
	uID, eID := userID(req), emailID(req)
	email, err := eh.repository.Find(uID, eID)
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}
	if email == nil {
		eh.errorResponse(w, ErrNotFound, http.StatusNotFound)
		return
	}
	if email.Status != "created" {
		eh.errorResponse(w, errors.New("email already sent"), http.StatusUnprocessableEntity)
		return
	}

	// TODO: schedule send
	log.Printf("Scheduling send for email %d\n", email.ID)

	email.Status = models.EmailQueued
	email, err = eh.repository.Update(uID, email)
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, email)

	return
}
