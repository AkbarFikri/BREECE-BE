package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/gocron"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/mailer"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type TicketService interface {
	ConfirmedPayment(invoiceId string) error
	FailurePayment(invoiceId string) error
}

type ticketService struct {
	EventRepository   repository.EventRepository
	InvoiceRepository repository.InvoiceRepository
	TicketRepository  repository.TicketRepository
	UserRepository    repository.UserRepository
	Mailer            mailer.EmailService
}

func NewTicketService(er repository.EventRepository, ir repository.InvoiceRepository,
	tr repository.TicketRepository, ur repository.UserRepository, m mailer.EmailService) TicketService {
	return &ticketService{
		EventRepository:   er,
		InvoiceRepository: ir,
		TicketRepository:  tr,
		UserRepository:    ur,
		Mailer:            m,
	}
}

// ConfirmedPayment implements TicketService.
func (s *ticketService) ConfirmedPayment(invoiceId string) error {
	invoice, _ := s.InvoiceRepository.FindById(invoiceId)
	user, _ := s.UserRepository.FindById(invoice.UserID)
	event, _ := s.EventRepository.FindById(invoice.EventID)

	invoice.Status = "success"

	if err := s.InvoiceRepository.Update(invoice); err != nil {
		return err
	}

	ticket := entity.Ticket{
		ID:        uuid.NewString(),
		UserID:    invoice.UserID,
		InvoiceID: invoice.ID,
		EventID:   invoice.EventID,
		CreatedAt: time.Now(),
	}

	if err := s.TicketRepository.Insert(ticket); err != nil {
		return err
	}

	if event.Tempat == "Online" {
		event.Tempat = event.Link
	}

	email := model.EmailNotification{
		Subject:    "Event Notification",
		Email:      user.Email,
		Name:       user.FullName,
		EventTitle: event.Title,
		EventStart: event.StartAt.String(),
		Venue:      event.Tempat,
	}

	// event.StartAt.Add(-(3 * time.Hour))

	if err := gocron.ScheduleSendNotification(event.StartAt.Add(-(3 * time.Hour)), s.Mailer, email); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

// FailurePayment implements TicketService.
func (s *ticketService) FailurePayment(invoiceId string) error {
	invoice, _ := s.InvoiceRepository.FindById(invoiceId)

	if err := s.EventRepository.UpdateFailurePayment(invoice.EventID); err != nil {
		return err
	}

	invoice.Status = "failure"

	if err := s.InvoiceRepository.Update(invoice); err != nil {
		return err
	}

	return nil
}
