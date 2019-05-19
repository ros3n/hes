package handlers

import (
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/parsers/emails"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/validators"
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

	email := &models.Email{Status: "created"}
	changeSet.ApplyChanges(email)

	email, err = eh.repository.Create(email)
	if err != nil {
		eh.errorResponse(w, ErrServerError, http.StatusInternalServerError)
		return
	}

	eh.jsonResponse(w, email)

	return
}
