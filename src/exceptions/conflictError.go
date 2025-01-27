package exceptions

import "fmt"

type ConflictError struct {
	Errors []FieldError
}

func (e ConflictError) Error() string {
	if len(e.Errors) == 0 {
		return "Registration failed."
	}

	var errMsg string
	for _, fe := range e.Errors {
		errMsg += fmt.Sprintf("%s %s.\n", fe.Field, fe.Tag)
	}
	return errMsg
}

func (e *ConflictError) AppendError(fieldName string, tag string) {
	e.Errors = append(e.Errors, FieldError{Field: fieldName, Tag: tag})
}

func (e ConflictError) FieldErrors() []FieldError {
	return e.Errors
}
