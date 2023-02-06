package email

import (
	"net/smtp"
)

func SendEmail(to []string, message []byte) error {
	from := "gofurovmurtazoxon@gmail.com"
	password := "password"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err

}
