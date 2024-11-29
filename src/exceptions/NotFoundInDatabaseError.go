package exceptions

import "fmt"

type NotFoundInDatabaseError struct {
	Err string
}

func NewNotFoundInDatabaseError() NotFoundInDatabaseError {
	return NotFoundInDatabaseError{
		Err: "Not in Db",
	}
}

func (e NotFoundInDatabaseError) Error() string {
	return fmt.Sprintf("Not found error: %s", e.Err)
}
