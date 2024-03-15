package service

import (
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
}

type paymentService struct {
	InvoiceRepository repository.InvoiceRepository
	EventRepository   repository.EventRepository
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
	event, err := s.EventRepository.FindForBooking(req.EventID)
	if err != nil {
		if err.Error() == "ticket is sold out" {
			return model.ServiceResponse{
				Code:    http.StatusUnprocessableEntity,
				Error:   true,
				Message: "Ticket is sold out.",
			}, err
		} else {
			return model.ServiceResponse{
				Code:    http.StatusInternalServerError,
				Error:   true,
				Message: "Something went wrong, failed to find event with id provided",
			}, err
		}
	}

	invoice := entity.Invoice{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		EventID:   event.ID,
		Amount:    int64(event.Price),
		Status:    "Pending",
		CreatedAt: time.Now(),
	}

	var snapResp *snap.Response

	if event.Price == 0 {
		invoice.Snap = "Success"
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
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return true
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				return false
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false
}
