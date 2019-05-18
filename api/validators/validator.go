package validators

type Validator interface {
	// Validate method performs data validation
	Validate()

	// Valid method informs if data validation has been successful
	Valid() bool
}

const cannotBeBlankError = "can't be blank"

// ValidationError contains a single error message related to a field
type ValidationError struct {
	Field   string
	Message string
}

// BaseValidator provides  basic utility functions and a container for error messages
type BaseValidator struct {
	Errors []ValidationError
}

func (v *BaseValidator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *BaseValidator) addError(field, message string) {
	v.Errors = append(v.Errors, ValidationError{field, message})
}
