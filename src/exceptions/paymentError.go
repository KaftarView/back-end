package exceptions

import "fmt"

type PaymentError struct {
	Message string
}

func NewPaymentError() PaymentError {
	return PaymentError{Message: "Unknown error occurred during payment processing"}
}

func (e PaymentError) Error() string {
	return e.Message
}

type InvalidPaymentRequestError struct {
	Field   string
	Details string
}

func (e InvalidPaymentRequestError) Error() string {
	return fmt.Sprintf("Invalid Payment Request - Field: %s, Details: %s", e.Field, e.Details)
}

type PaymentServerError struct {
	Message string
}

func NewPaymentServerError() PaymentServerError {
	return PaymentServerError{Message: "Failed to initialize payment service"}
}

func (e PaymentServerError) Error() string {
	return fmt.Sprintf("Payment Server Error - Message: %s", e.Message)
}
