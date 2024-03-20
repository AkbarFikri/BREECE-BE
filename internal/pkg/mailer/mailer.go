package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

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
	template *template.Template
}

func NewMailer() EmailService {
	username := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	t := template.Must(template.ParseGlob("internal/pkg/mailer/template/*.html"))
	return &mailer{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		template: t,
	}
}

func (m *mailer) SendOTP(data model.EmailOTP) error {
	buffer := new(bytes.Buffer)

	if err := m.template.ExecuteTemplate(buffer, "otp_email.html", data); err != nil {
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: Text/html; charset=\"iso-8859-1\";\n\n"
	fromUser := fmt.Sprintf("From: BREECE <%s>\n", m.Username)
	toUser := fmt.Sprintf("To: %s\n", data.Email)
	subjectEmail := fmt.Sprintf("Subject: %s\n", data.Subject)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	body := []byte(fromUser + toUser + subjectEmail + mime + buffer.String())

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.Username, []string{data.Email}, body)
}

func (m *mailer) SendNotification(data model.EmailNotification) error {
	buffer := new(bytes.Buffer)

	if err := m.template.ExecuteTemplate(buffer, "notification_event_email.html", data); err != nil {
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: Text/html; charset=\"iso-8859-1\";\n\n"
	fromUser := fmt.Sprintf("From: BREECE <%s>\n", m.Username)
	toUser := fmt.Sprintf("To: %s\n", data.Email)
	subjectEmail := fmt.Sprintf("Subject: %s\n", data.Subject)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	body := []byte(fromUser + toUser + subjectEmail + mime + buffer.String())

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.Username, []string{data.Email}, body)
}

func (m *mailer) SendInvoice() error {
	return nil
}
