package exceptions

type RateLimitError struct {
	Message string
}

func NewRateLimitError() RateLimitError {
	return RateLimitError{
		Message: "Rate limit exceeded. Please try again later.",
	}
}

func (e RateLimitError) Error() string {
	return e.Message
}
