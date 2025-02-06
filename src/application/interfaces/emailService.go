package application_interfaces

type EmailService interface {
	SendEmail(toEmail string, subject string, templateFile string, data interface{})
}
