package repositories

import (
	"github.com/ros3n/hes/api/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	userID        = "user-id"
	emailID int64 = 1
)

type SimpleEmailsRepoTestSuite struct {
	suite.Suite
	repo *SimpleEmailsRepository
}

func (suite *SimpleEmailsRepoTestSuite) SetupTest() {
	suite.repo = NewSimpleEmailsRepository()
}

func (suite *SimpleEmailsRepoTestSuite) TestFind() {
	email := testEmail(emailID)
	suite.repo.emails[userID] = map[int64]*models.Email{emailID: email}

	fetchedEmail, _ := suite.repo.Find(userID, emailID)

	suite.True(email != fetchedEmail)
	suite.Equal(email, fetchedEmail)

	nonExistentEmail, _ := suite.repo.Find(userID, 0)
	suite.Nil(nonExistentEmail)
}

func (suite *SimpleEmailsRepoTestSuite) TestCreate() {
	email := testEmail(0)

	email, _ = suite.repo.Create(userID, email)
	suite.Equal(emailID, email.ID)
	suite.Equal(int64(2), suite.repo.autoincrementId)
	suite.Equal(userID, email.UserID)
}

func (suite *SimpleEmailsRepoTestSuite) TestAll() {
	email := testEmail(emailID)
	suite.repo.emails[userID] = map[int64]*models.Email{emailID: email}

	var otherEmailID int64 = emailID + 1
	otherEmail := testEmail(otherEmailID)
	otherUserID := "other-user"
	suite.repo.emails[otherUserID] = map[int64]*models.Email{otherEmailID: otherEmail}

	emails, _ := suite.repo.All(userID)
	suite.Equal(1, len(emails))
	suite.Equal(email, emails[0])
}

func TestSimpleEmailsRepository(t *testing.T) {
	suite.Run(t, new(SimpleEmailsRepoTestSuite))
}

func testEmail(id int64) *models.Email {
	return &models.Email{
		ID:         id,
		Sender:     "sender@example.com",
		Recipients: []string{"recipient@example.com"},
		Subject:    "subject",
		Message:    "message",
		Status:     models.EmailCreated,
	}
}
