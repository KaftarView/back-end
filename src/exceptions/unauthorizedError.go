package exceptions

import "fmt"

type UnauthorizedError struct {
	Err string
}

func NewUnauthorizedError() UnauthorizedError {
	return UnauthorizedError{
		Err: "There is an issue with your authentication",
	}
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized error: %s", e.Err)
}
