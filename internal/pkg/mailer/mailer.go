package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type EmailService interface {
	SendOTP(data model.EmailOTP) error
	SendNotification(data model.EmailNotification) error
}

type mailer struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewMailer() EmailService {
	username := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	return &mailer{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (m *mailer) SendOTP(data model.EmailOTP) error {
	send := strings.ReplaceAll(otp_email, "{{.otp}}", data.Otp)
	send = strings.ReplaceAll(send, "{{.name}}", data.Name)

	mime := "MIME-version: 1.0;\nContent-Type: Text/html; charset=\"iso-8859-1\";\n\n"
	fromUser := fmt.Sprintf("From: BREECE <%s>\n", m.Username)
	toUser := fmt.Sprintf("To: %s\n", data.Email)
	subjectEmail := fmt.Sprintf("Subject: %s\n", data.Subject)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	body := []byte(fromUser + toUser + subjectEmail + mime + send)

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.Username, []string{data.Email}, body)
}

func (m *mailer) SendNotification(data model.EmailNotification) error {
	send := strings.ReplaceAll(send_notification, "{{.event-title}}", data.EventTitle)
	send = strings.ReplaceAll(send, "{{.name}}", data.Name)
	send = strings.ReplaceAll(send, "{{.event-start}}", data.EventStart)
	send = strings.ReplaceAll(send, "{{.venue}}", data.Venue)

	mime := "MIME-version: 1.0;\nContent-Type: Text/html; charset=\"iso-8859-1\";\n\n"
	fromUser := fmt.Sprintf("From: BREECE <%s>\n", m.Username)
	toUser := fmt.Sprintf("To: %s\n", data.Email)
	subjectEmail := fmt.Sprintf("Subject: %s\n", data.Subject)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	body := []byte(fromUser + toUser + subjectEmail + mime + send)

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.Username, []string{data.Email}, body)
}

func (m *mailer) SendInvoice() error {
	return nil
}
