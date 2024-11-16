package application

import (
	"first-project/src/bootstrap"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	emailInfo *bootstrap.EmailInfo
}

func NewEmailService(emailInfo *bootstrap.EmailInfo) *EmailService {
	return &EmailService{
		emailInfo: emailInfo,
	}
}

// TODO: change this to communication and then email sender!
func (emailService *EmailService) SendVerificationEmail(username string, toEmail string, otp string) {
	from := emailService.emailInfo.EmailFrom
	password := emailService.emailInfo.EmailPassword
	smtpHost := emailService.emailInfo.SMTPHost
	smtpPort := emailService.emailInfo.SMTPPort
	verificationLink := "http://localhost:8080/v1/register/activate"
	// TODO: translation and template
	message := []byte(
		fmt.Sprintf(
			"dear %s;\nPlease visit: %s and input otp below to activate your account:\n%s",
			username, verificationLink, otp),
	)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		panic(err)
	}
}
