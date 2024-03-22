package service

import (
	"errors"
	"fmt"
	"net/http"
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
	FetchUserTicketHistory(user model.UserTokenData) (model.ServiceResponse, error)
	FetchParticipantTicket(user model.UserTokenData, eventId string) (model.ServiceResponse, error)
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

	// TODO Hindari compare string langsung.
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

func (s *ticketService) FetchUserTicketHistory(user model.UserTokenData) (model.ServiceResponse, error) {
	tickets, err := s.TicketRepository.FindByUserId(user.ID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find ticket history",
		}, err
	}

	if len(tickets) == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Record not found",
		}, errors.New("record not found")
	}

	var res []model.TicketUserResponse

	for _, t := range tickets {
		dumpEvent := model.EventResponse{
			ID:           t.Event.ID,
			CategoryID:   t.Event.CategoryID,
			Title:        t.Event.Title,
			Description:  t.Event.Description,
			Place:        t.Event.Tempat,
			Speakers:     t.Event.Speakers,
			SpeakersRole: t.Event.SpeakersRole,
			Date:         t.Event.Date.String(),
			StartAt:      t.Event.StartAt.String(),
			Link:         t.Event.Link,
			Price:        t.Event.Price,
			OrganizeBy:   t.Event.OrganizeBy,
			IsPublic:     t.Event.IsPublic,
		}

		dump := model.TicketUserResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			EventID:   t.EventID,
			InvoiceID: t.InvoiceID,
			CreatedAt: t.CreatedAt,
			Event:     dumpEvent,
		}

		res = append(res, dump)
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all tickets",
		Data:    res,
	}, nil
}

func (s *ticketService) FetchParticipantTicket(user model.UserTokenData, eventId string) (model.ServiceResponse, error) {
	event, err := s.EventRepository.FindById(eventId)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Something went wrong, failed to find event",
		}, err
	}

	if event.OrganizeBy != user.ID {
		return model.ServiceResponse{
			Code:    http.StatusForbidden,
			Error:   true,
			Message: "You're not allowed to look at this data",
		}, err
	}

	tickets, err := s.TicketRepository.FindByEventId(eventId)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find participants history",
		}, err
	}

	if len(tickets) == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Record not found",
		}, errors.New("record not found")
	}

	var res []model.TicketOrganizerResponse

	for _, t := range tickets {
		dumpUser := model.ProfileUserResponse{
			ID:          t.ID,
			Email:       t.User.Email,
			FullName:    t.User.FullName,
			NimNik:      t.User.NimNik,
			Prodi:       t.User.Prodi,
			Universitas: t.User.Universitas,
		}

		dump := model.TicketOrganizerResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			EventID:   t.EventID,
			InvoiceID: t.InvoiceID,
			CreatedAt: t.CreatedAt,
			User:      dumpUser,
		}

		res = append(res, dump)
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all participants",
		Data:    res,
	}, nil
}
