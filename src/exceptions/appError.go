package exceptions

import "fmt"

type AppError struct {
	Field string
	Tag   string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error 400 Bad Request: %s, %s", e.Field, e.Tag)
}

func NewAppError(field, tag string) *AppError {
	return &AppError{
		Field: field,
		Tag:   tag,
	}
}
