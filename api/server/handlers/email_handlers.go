package handlers

import (
	"errors"
	"github.com/ros3n/hes/api/parsers/emails"
	"github.com/ros3n/hes/api/services"
	"github.com/ros3n/hes/api/validators"
	"net/http"
)

type EmailsAPIHandler struct {
	*apiHandler
	emailService *services.EmailService
}

func NewEmailsAPIHandler(service *services.EmailService) *EmailsAPIHandler {
	return &EmailsAPIHandler{apiHandler: &apiHandler{}, emailService: service}
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

	email, err := eh.emailService.Create(userID(req), parser.Data())

	if err == nil {
		eh.jsonResponseWithStatus(w, email, http.StatusCreated)
		return
	}

	if evErr, ok := err.(*validators.EmailValidator); ok {
		eh.jsonResponseWithStatus(w, evErr.Errors(), http.StatusUnprocessableEntity)
		return
	}

	eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
	return
}

func (eh *EmailsAPIHandler) ListEmails(w http.ResponseWriter, req *http.Request) {
	allEmails, err := eh.emailService.All(userID(req))
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, allEmails)
	return
}

func (eh *EmailsAPIHandler) SendEmail(w http.ResponseWriter, req *http.Request) {
	uID, eID := userID(req), emailID(req)
	email, err := eh.emailService.Find(uID, eID)
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

	err = eh.emailService.ScheduleSend(req.Context(), email)
	if err != nil {
		eh.errorResponse(w, errors.New("failed to schedule send"), http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, email)
	return
}
