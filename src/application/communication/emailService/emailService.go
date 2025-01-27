package application_communication

import (
	"bytes"
	"first-project/src/bootstrap"
	"html/template"
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

func (emailService *EmailService) SendEmail(toEmail string, subject string, templateFile string, data interface{}) {
	from := emailService.emailInfo.EmailFrom
	password := emailService.emailInfo.EmailPassword
	smtpHost := emailService.emailInfo.SMTPHost
	smtpPort := emailService.emailInfo.SMTPPort

	tmpl, err := template.ParseFiles("src/application/communication/emailService/templates/" + templateFile)
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer
	body.Write([]byte("To: " + toEmail + "\r\n"))
	body.Write([]byte("Subject: " + subject + "\r\n"))
	body.Write([]byte("MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"))
	err = tmpl.Execute(&body, data)
	if err != nil {
		panic(err)
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, body.Bytes())
	if err != nil {
		panic(err)
	}
}
