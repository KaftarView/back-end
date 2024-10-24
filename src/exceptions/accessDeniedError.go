package exceptions

import "fmt"

type AuthError struct {
	Err string
}

func NewAuthError() AuthError {
	return AuthError{
		Err: "Access denied",
	}
}

func (e AuthError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Err)
}
