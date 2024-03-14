package service

import (
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type PaymentService interface {
	GenerateUrlAndToken(user model.UserTokenData, req model.PaymentRequest) (model.ServiceResponse, error)
	VerifyPayment(data map[string]interface{}) error
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

	// 3. Request create Snap transaction to Midtrans
	snapResp, _ := s.Client.CreateTransaction(payReq)
	// if err != nil {
	// 	return model.ServiceResponse{
	// 		Code:    http.StatusInternalServerError,
	// 		Error:   true,
	// 		Message: "Something went wrong, failed to create payment snap and token.",
	// 	}, err
	// }

	invoice.Snap = snapResp.RedirectURL

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
func (*paymentService) VerifyPayment(data map[string]interface{}) error {
	panic("unimplemented")
}
