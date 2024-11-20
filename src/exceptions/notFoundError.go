package exceptions

import "fmt"

type NotFoundError struct {
	ErrorField string
}

func (e NotFoundError) Error() string {
	if len(e.ErrorField) == 0 {
		return "404 Page Not Found."
	}
	return fmt.Sprintf("404 %s Not Found", e.ErrorField)
}

func (e NotFoundError) FieldErrors() string {
	return e.ErrorField
}
