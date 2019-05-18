package utils

// StrPointer returns a pointer to a given string. It is useful
// when one have to differentiate between nil and empty strings.
func StrPointer(str string) *string {
	return &str
}
