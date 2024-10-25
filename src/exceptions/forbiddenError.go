package exceptions

import "fmt"

type ForbiddenError struct {
	Err string
}

func NewForbiddenError() ForbiddenError {
	return ForbiddenError{
		Err: "Access denied",
	}
}

func (e ForbiddenError) Error() string {
	return fmt.Sprintf("Forbidden error: %s", e.Err)
}
