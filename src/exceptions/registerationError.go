package exceptions

import (
	"fmt"
)

type FieldError struct {
	Field string
	Tag   string
}

type UserRegistrationError struct {
	Errors []FieldError
}

func (e UserRegistrationError) Error() string {
	if len(e.Errors) == 0 {
		return "Registration failed."
	}

	var errMsg string
	for _, fe := range e.Errors {
		errMsg += fmt.Sprintf("%s %s.\n", fe.Field, fe.Tag)
	}
	return errMsg
}

func (e *UserRegistrationError) AppendError(fieldName string, tag string) {
	e.Errors = append(e.Errors, FieldError{Field: fieldName, Tag: tag})
}

func (e UserRegistrationError) FieldErrors() []FieldError {
	return e.Errors
}
