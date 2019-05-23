package server

import (
	"context"
	"github.com/ros3n/hes/api/models"
	"github.com/ros3n/hes/api/repositories"
	"github.com/ros3n/hes/api/server/handlers"
	"github.com/ros3n/hes/api/server/middleware"
	"github.com/ros3n/hes/api/services"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EmailsAPITestSuite struct {
	suite.Suite
	authenticate bool
	handler      http.Handler
}

const (
	user     = "test"
	password = "test"
	userID   = "1"
)

func (suite *EmailsAPITestSuite) SetupSuite() {
	authService := middleware.NewBasicAuthenticator(user, password, userID)
	repo := repositories.NewSimpleEmailsRepository()
	msgSender := &FakeMessageSender{}
	emailService := services.NewEmailService(repo, msgSender)
	emailsHandler := handlers.NewEmailsAPIHandler(emailService)

	suite.handler = newRouter(authService, emailsHandler)
}

func (suite *EmailsAPITestSuite) SetupTest() {
	suite.AuthenticateRequests(true)
}

func (suite *EmailsAPITestSuite) AuthenticateRequests(authenticate bool) {
	suite.authenticate = authenticate
}

func (suite *EmailsAPITestSuite) TestAPIAuthentication() {
	suite.AuthenticateRequests(false)
	resp := suite.makeRequest("POST", "/emails", "")

	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (suite *EmailsAPITestSuite) TestEmailValidation() {
	testCases := []struct {
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			`{sender: id}`,
			http.StatusBadRequest,
			`{"error":"malformed request body"}`,
		}, {
			`{"sender":"","recipients": ["test@example.com"],"subject":"test","message":"test"}`,
			http.StatusUnprocessableEntity,
			`{"errors":[{"field":"sender","message":"can't be blank"}]}`,
		}, {
			`{"sender":"test","recipients": ["test@example.com"],"subject":"test","message":"test"}`,
			http.StatusUnprocessableEntity,
			`{"errors":[{"field":"sender","message":"test is not a valid email address"}]}`,
		}, {
			`{"sender":"test@example.com","subject":"test","message":"test"}`,
			http.StatusUnprocessableEntity,
			`{"errors":[{"field":"recipients","message":"can't be blank"}]}`,
		}, {
			`{"recipients": ["valid@example.com", "invalid"],"sender":"test@example.com","subject":"test","message":"test"}`,
			http.StatusUnprocessableEntity,
			`{"errors":[{"field":"recipients","message":"invalid is not a valid email address"}]}`,
		},
	}

	for _, testCase := range testCases {
		resp := suite.makeRequest("POST", "/emails", testCase.requestBody)

		suite.Equal(testCase.expectedStatus, resp.StatusCode)
		suite.Equal(testCase.expectedBody, suite.getResponseBody(resp))
	}
}

func (suite *EmailsAPITestSuite) TestEmailCreation() {
	testCases := []struct {
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			`{"recipients": ["valid@example.com"],"sender":"test@example.com","subject":"test","message":"test"}`,
			http.StatusCreated,
			`{"id":1,"user_id":"1","sender":"test@example.com","recipients":["valid@example.com"],"subject":"test","message":"test","status":"created"}`,
		},
	}

	for _, testCase := range testCases {
		resp := suite.makeRequest("POST", "/emails", testCase.requestBody)

		suite.Equal(testCase.expectedStatus, resp.StatusCode)
		suite.Equal(testCase.expectedBody, suite.getResponseBody(resp))
	}
}

func (suite *EmailsAPITestSuite) makeRequest(method, path, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		suite.FailNow(err.Error())
	}

	if suite.authenticate {
		req.SetBasicAuth(user, password)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.handler.ServeHTTP(rr, req)

	return rr.Result()
}

func (suite *EmailsAPITestSuite) getResponseBody(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}
	return string(body)
}

func TestEmailsAPI(t *testing.T) {
	suite.Run(t, new(EmailsAPITestSuite))
}

type FakeMessageSender struct{}

func (fms *FakeMessageSender) SendEmail(ctx context.Context, email *models.Email) error {
	return nil
}
