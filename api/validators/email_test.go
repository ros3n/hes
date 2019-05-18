package validators

import (
	"github.com/ros3n/hes/lib/utils"
	"testing"

	"github.com/ros3n/hes/api/models"
	"github.com/stretchr/testify/suite"
)

type EmailValidatorTestSuite struct {
	suite.Suite
}

func testChangeSet() *models.EmailChangeSet {
	return &models.EmailChangeSet{
		Sender:     utils.StrPointer("sender@example.com"),
		Recipients: []string{"recipient@example.com"},
		Subject:    utils.StrPointer("test subject"),
		Message:    utils.StrPointer("test message"),
	}
}

func (suite *EmailValidatorTestSuite) runTestCase(changeSet *models.EmailChangeSet, key string, valid bool, errorMsg string) {
	validator := NewEmailValidator(changeSet)
	validator.Validate()

	suite.Equal(valid, validator.Valid())
	var errors []ValidationError
	if errorMsg != "" {
		errors = append(errors, ValidationError{key, errorMsg})
	}
	suite.Equal(errors, validator.Errors)
}

func (suite *EmailValidatorTestSuite) TestSenderValidation() {
	testCases := []struct {
		sender *string
		valid  bool
		error  string
	}{
		{
			nil, false, cannotBeBlankError,
		},
		{
			utils.StrPointer(""), false, cannotBeBlankError,
		},
		{
			utils.StrPointer("invalidemail"), false, "invalidemail is not a valid email address",
		},
		{
			utils.StrPointer("email@example.com"), true, "",
		},
	}

	for _, testCase := range testCases {
		changeSet := testChangeSet()
		changeSet.Sender = testCase.sender
		suite.runTestCase(changeSet, senderKey, testCase.valid, testCase.error)
	}
}
func (suite *EmailValidatorTestSuite) TestRecipientsValidation() {
	testCases := []struct {
		recipients []string
		valid      bool
		error      string
	}{
		{
			nil, false, cannotBeBlankError,
		},
		{
			[]string{}, false, cannotBeBlankError,
		},
		{
			[]string{"invalidemail"}, false, "invalidemail is not a valid email address",
		},
		{
			[]string{"email@example.com"}, true, "",
		},
	}

	for _, testCase := range testCases {
		changeSet := testChangeSet()
		changeSet.Recipients = testCase.recipients
		suite.runTestCase(changeSet, recipientsKey, testCase.valid, testCase.error)
	}
}

func (suite *EmailValidatorTestSuite) TestSubjectValidation() {
	testCases := []struct {
		subject *string
		valid   bool
		error   string
	}{
		{
			nil, false, cannotBeBlankError,
		},
		{
			utils.StrPointer(""), false, cannotBeBlankError,
		},
		{
			utils.StrPointer("test subject"), true, "",
		},
	}

	for _, testCase := range testCases {
		changeSet := testChangeSet()
		changeSet.Subject = testCase.subject
		suite.runTestCase(changeSet, subjectKey, testCase.valid, testCase.error)
	}
}

func (suite *EmailValidatorTestSuite) TestMessageValidation() {
	testCases := []struct {
		message *string
		valid   bool
		error   string
	}{
		{
			nil, false, cannotBeBlankError,
		},
		{
			utils.StrPointer(""), false, cannotBeBlankError,
		},
		{
			utils.StrPointer("test message"), true, "",
		},
	}

	for _, testCase := range testCases {
		changeSet := testChangeSet()
		changeSet.Message = testCase.message
		suite.runTestCase(changeSet, messageKey, testCase.valid, testCase.error)
	}
}

func TestEmailValidator(t *testing.T) {
	suite.Run(t, new(EmailValidatorTestSuite))
}
