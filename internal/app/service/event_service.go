package service

import (
	"errors"
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
	FetchEvent(user model.UserTokenData, params model.FilterParam) (model.ServiceResponse, error)
	FetchEventDetails(id string) (model.ServiceResponse, error)
	FetchCategory() (model.ServiceResponse, error)
}

type eventService struct {
	EventRepository    repository.EventRepository
	CategoryRepository repository.CategoryRepository
	SupabaseBucket     *supabasestorageuploader.Client
}

func NewEventService(er repository.EventRepository,
	sb *supabasestorageuploader.Client,
	cr repository.CategoryRepository) EventService {
	return &eventService{
		EventRepository:    er,
		SupabaseBucket:     sb,
		CategoryRepository: cr,
	}
}

func (s *eventService) PostEvent(user model.UserTokenData, req model.EventRequest) (model.ServiceResponse, error) {
	date, err := time.Parse("2006-01-02 15:04:05 -0700 MST", req.Date)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid time format for field date",
		}, err
	}

	startAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", req.StartAt)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid time format for field start_at",
		}, err
	}

	datenow, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", time.Now().UTC().Format("2006-01-02")+" 00:00:00 +0000 UTC")
	timenow, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", time.Now().UTC().Format("2006-01-02 15:04:05 -0700 MST"))
	if date.After(datenow) || startAt.After(timenow) {
		return model.ServiceResponse{
			Code:    http.StatusForbidden,
			Error:   true,
			Message: "Invalid time request, the event holding time cannot be less than the current time",
		}, err
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
		// TODO Upload the default banner and replace the string "default"
		event.BannerUrl = "default"
	} else {
		req.Banner.Filename = "banner_" + event.ID + ".png"
		bannerUrl, err := s.SupabaseBucket.Upload(req.Banner)
		if err != nil {
			return model.ServiceResponse{
				Code:    http.StatusInternalServerError,
				Error:   true,
				Message: "Failed to upload banner to bucket",
			}, err
		}
		event.BannerUrl = bannerUrl
	}

	if err := s.EventRepository.Insert(event); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went frong, failed to create event.",
		}, err
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

func (s *eventService) FetchEvent(user model.UserTokenData, params model.FilterParam) (model.ServiceResponse, error) {
	if params.Page == 0 {
		params.Page = 1
	}

	params.IsPublic = !user.IsBrawijaya

	if params.Date != "" {
		_, err := time.Parse("2006-01-02", params.Date)
		if err != nil {
			return model.ServiceResponse{
				Code:    http.StatusInternalServerError,
				Error:   true,
				Message: "Invalid time format on field date query param",
			}, err
		}
	}

	events, err := s.EventRepository.FindWithFilter(params)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Something went wrong, events with filter params provided is not found.",
		}, err
	}

	// TODO Errornya perbaiki!
	if len(events) == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Events with filter params provided is not found.",
		}, errors.New("Record not found")
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all events",
		Data:    events,
	}, nil
}

func (s *eventService) FetchEventDetails(id string) (model.ServiceResponse, error) {
	event, err := s.EventRepository.FindById(id)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Event with id provided is not found",
		}, err
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully found event",
		Data:    event,
	}, nil
}

func (s *eventService) FetchCategory() (model.ServiceResponse, error) {
	categories, err := s.CategoryRepository.FindAll()
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find category",
		}, err
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all categories",
		Data:    categories,
	}, nil
}
