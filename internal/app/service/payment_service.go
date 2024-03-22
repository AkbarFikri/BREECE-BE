package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type PaymentService interface {
	GenerateUrlAndToken(user model.UserTokenData, req model.PaymentRequest) (model.ServiceResponse, error)
	VerifyPayment(orderId string) bool
	FetchPaymentHistory(user model.UserTokenData) (model.ServiceResponse, error)
}

type paymentService struct {
	InvoiceRepository repository.InvoiceRepository
	EventRepository   repository.EventRepository
	TicketRepository  repository.TicketRepository
	Client            snap.Client
}

func NewPaymentService(ir repository.InvoiceRepository, er repository.EventRepository) PaymentService {
	var client snap.Client
	env := midtrans.Sandbox
	client.New(os.Getenv("MIDTRANS_KEY"), env)

	return &paymentService{
		Client:            client,
		InvoiceRepository: ir,
		EventRepository:   er,
	}
}

// GenerateUrlAndToken implements PaymentService.
func (s *paymentService) GenerateUrlAndToken(user model.UserTokenData, req model.PaymentRequest) (model.ServiceResponse, error) {
	event, err := s.EventRepository.FindById(req.EventID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Something went wrong, failed to find event",
		}, err
	}

	if event.TicketQty == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotAcceptable,
			Error:   true,
			Message: "Ticket is sold out",
		}, err
	}

	if err := s.EventRepository.UpdateTicketDecrement(event); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Something went wrong, failed to update event ticket_qty",
		}, err
	}

	invoice := entity.Invoice{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		EventID:   event.ID,
		Amount:    int64(event.Price),
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	var snapResp *snap.Response

	if event.Price == 0 {
		invoice.Snap = "success"
		invoice.Status = "success"

		ticket := entity.Ticket{
			ID:        uuid.NewString(),
			UserID:    user.ID,
			InvoiceID: invoice.ID,
			EventID:   event.ID,
			CreatedAt: time.Now().UTC(),
		}

		if err := s.TicketRepository.Insert(ticket); err != nil {
			fmt.Print("Error")
		}
	} else {
		payReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  invoice.ID,
				GrossAmt: invoice.Amount + 2000,
			},
			Expiry: &snap.ExpiryDetails{
				Duration: 15,
				Unit:     "minute",
			},
		}

		snapResp, _ = s.Client.CreateTransaction(payReq)

		invoice.Snap = snapResp.RedirectURL
	}

	if err := s.InvoiceRepository.Insert(invoice); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to create invoice.",
		}, err
	}

	res := model.PaymentResponse{
		SnapUrl: snapResp.RedirectURL,
		Token:   snapResp.Token,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully create payment method",
		Data:    res,
	}, nil
}

// VerifyPayment implements PaymentService.
func (s *paymentService) VerifyPayment(orderId string) bool {
	var client coreapi.Client
	env := midtrans.Sandbox
	client.New(os.Getenv("MIDTRANS_KEY"), env)

	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return false
	} else {
		if transactionStatusResp != nil {
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					return false
				} else if transactionStatusResp.FraudStatus == "accept" {
					return true
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return true
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				return false
			}
		}
	}
	return false
}

func (s *paymentService) FetchPaymentHistory(user model.UserTokenData) (model.ServiceResponse, error) {
	invoices, err := s.InvoiceRepository.FindByUserId(user.ID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find payment history",
		}, err
	}

	if len(invoices) == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Record not found",
		}, errors.New("record not found")
	}

	var res []model.PaymentHistoryResponse

	for _, invoice := range invoices {
		dumpEvent := model.EventResponse{
			ID:           invoice.Event.ID,
			CategoryID:   invoice.Event.CategoryID,
			Title:        invoice.Event.Title,
			Description:  invoice.Event.Description,
			Place:        invoice.Event.Tempat,
			Speakers:     invoice.Event.Speakers,
			SpeakersRole: invoice.Event.SpeakersRole,
			Date:         invoice.Event.Date.String(),
			StartAt:      invoice.Event.StartAt.String(),
			Link:         invoice.Event.Link,
			Price:        invoice.Event.Price,
			OrganizeBy:   invoice.Event.OrganizeBy,
			IsPublic:     invoice.Event.IsPublic,
		}

		dump := model.PaymentHistoryResponse{
			ID:     invoice.ID,
			Amount: invoice.Amount,
			Status: invoice.Status,
			Event:  dumpEvent,
		}

		res = append(res, dump)
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all payment history",
		Data:    res,
	}, nil
}
