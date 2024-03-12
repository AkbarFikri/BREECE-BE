package service

import (
	"net/http"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/google/uuid"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type EventService interface {
	PostEvent(user model.UserTokenData, req model.EventRequest) (model.ServiceResponse, error)
}

type eventService struct {
	EventRepository repository.EventRepository
	SupabaseBucket  *supabasestorageuploader.Client
}

func NewEventService(er repository.EventRepository,
	sb *supabasestorageuploader.Client) EventService {
	return &eventService{
		EventRepository: er,
	}
}

func (s *eventService) PostEvent(user model.UserTokenData, req model.EventRequest) (model.ServiceResponse, error) {
	date, err := time.Parse("2006-01-02 15:04:05 -0700 MST", req.Date)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid time format for field date",
		}, nil
	}

	startAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", req.StartAt)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid time format for field start_at",
		}, nil
	}

	event := entity.Event{
		ID:           uuid.NewString(),
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Description:  req.Description,
		Tempat:       req.Place,
		Speakers:     req.Speakers,
		SpeakersRole: req.SpeakersRole,
		Date:         date,
		StartAt:      startAt,
		Link:         req.Link,
		Price:        req.Price,
		TicketQty:    req.TicketQty,
		OrganizeBy:   user.ID,
		IsPublic:     req.IsPublic,
	}

	if req.Banner == nil {
		event.BannerUrl = "default"
	} else {
		req.Banner.Filename = "banner_" + event.ID + ".png"
		bannerUrl, err := s.SupabaseBucket.Upload(req.Banner)
		if err != nil {
			return model.ServiceResponse{
				Code:    http.StatusInternalServerError,
				Error:   true,
				Message: "Failed to upload banner to bucket",
			}, nil
		}
		event.BannerUrl = bannerUrl
	}

	if err := s.EventRepository.Insert(event); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went frong, failed to create event.",
		}, nil
	}

	res := model.EventResponse{
		ID:           event.ID,
		CategoryID:   event.CategoryID,
		Title:        event.Title,
		Description:  event.Description,
		Place:        event.Tempat,
		Speakers:     event.Speakers,
		SpeakersRole: event.SpeakersRole,
		BannerUrl:    event.BannerUrl,
		Date:         event.Date.String(),
		StartAt:      event.StartAt.String(),
		Link:         event.Link,
		Price:        event.Price,
		TicketQty:    event.TicketQty,
		OrganizeBy:   event.OrganizeBy,
		IsPublic:     event.IsPublic,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully to create event.",
		Data:    res,
	}, nil
}
