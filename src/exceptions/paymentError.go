package exceptions

import "fmt"

type PaymentError struct {
	Message string
	Code    int
}

func (e PaymentError) Error() string {
	return fmt.Sprintf("Payment Error - Code: %d, Message: %s", e.Code, e.Message)
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

func (e PaymentServerError) Error() string {
	return fmt.Sprintf("Payment Server Error - Message: %s", e.Message)
}

type UnhandledPaymentError struct {
	Message string
}

func (e UnhandledPaymentError) Error() string {
	return fmt.Sprintf("Erro in Payment %s", e.Message)
}
