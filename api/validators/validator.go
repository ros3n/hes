package validators

type Validator interface {
	// Validate method performs data validation
	Validate()

	// Valid method informs if data validation has been successful
	Valid() bool

	// errors returns an array of validation errors
	Errors() []ValidationError
	Error() string
}

const cannotBeBlankError = "can't be blank"

// ValidationError contains a single error message related to a field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// BaseValidator provides  basic utility functions and a container for error messages
type BaseValidator struct {
	errors []ValidationError
}

func (v *BaseValidator) Valid() bool {
	return len(v.errors) == 0
}

func (v *BaseValidator) Errors() []ValidationError {
	return v.errors
}

func (v *BaseValidator) APIErrors() map[string][]ValidationError {
	return map[string][]ValidationError{"errors": v.errors}
}

func (v *BaseValidator) Error() string {
	return "record is not valid"
}

func (v *BaseValidator) addError(field, message string) {
	v.errors = append(v.errors, ValidationError{field, message})
}
