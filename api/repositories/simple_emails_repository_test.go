package repositories

import (
	"github.com/ros3n/hes/api/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

const emailID = 1

type SimpleEmailsRepoTestSuite struct {
	suite.Suite
	repo *SimpleEmailsRepository
}

func (suite *SimpleEmailsRepoTestSuite) SetupTest() {
	suite.repo = NewSimpleEmailsRepository()
}

func (suite *SimpleEmailsRepoTestSuite) TestFind() {
	email := testEmail(emailID)
	suite.repo.emails[emailID] = email

	fetchedEmail, _ := suite.repo.Find(emailID)

	suite.True(email != fetchedEmail)
	suite.Equal(email, fetchedEmail)

	nonExistentEmail, _ := suite.repo.Find(0)
	suite.Nil(nonExistentEmail)
}

func (suite *SimpleEmailsRepoTestSuite) TestCreate() {
	email := testEmail(0)

	email, _ = suite.repo.Create(email)
	suite.Equal(1, email.ID)
	suite.Equal(2, suite.repo.autoincrementId)
}

func (suite *SimpleEmailsRepoTestSuite) TestAll() {
	email := testEmail(emailID)
	suite.repo.emails[emailID] = email

	emails, _ := suite.repo.All()
	suite.Equal(1, len(emails))
	suite.Equal(email, emails[0])
}

func TestSimpleEmailsRepository(t *testing.T) {
	suite.Run(t, new(SimpleEmailsRepoTestSuite))
}

func testEmail(id int) *models.Email {
	return &models.Email{
		ID:         id,
		Sender:     "sender@example.com",
		Recipients: []string{"recipient@example.com"},
		Subject:    "subject",
		Message:    "message",
		Status:     models.EmailCreated,
	}
}
