package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func Send(to, subject, otp string, name string) error {
	username := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")

	send := strings.ReplaceAll(otp_email, "{{.otp}}", otp)
	send = strings.ReplaceAll(send, "{{.name}}", name)

	mime := "MIME-version: 1.0;\nContent-Type: Text/html; charset=\"iso-8859-1\";\n\n"
	fromUser := fmt.Sprintf("From: BREECE <%s>\n", username)
	toUser := fmt.Sprintf("To: %s\n", to)
	subjectEmail := fmt.Sprintf("Subject: %s\n", subject)

	auth := smtp.PlainAuth("", username, password, host)
	body := []byte(fromUser + toUser + subjectEmail + mime + send)

	return smtp.SendMail(host+":"+port, auth, username, []string{to}, body)
}
