package exceptions

import "fmt"

type LoginError struct {
	Err string
}

func NewLoginError() LoginError {
	return LoginError{
		Err: "Invalid username or password.",
	}
}

func (e LoginError) Error() string {
	return fmt.Sprintf("Login failed: %s", e.Err)
}
