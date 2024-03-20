package service

import (
	"net/http"

	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type AdminService interface {
	FetchOrganizer() (model.ServiceResponse, error)
	FetchOrganizerDetail(id string) (model.ServiceResponse, error)
	FetchUser() (model.ServiceResponse, error)
	VerifyOrganizer(req model.OrganizerVerifyRequest) (model.ServiceResponse, error)
}

type adminService struct {
	UserRepository repository.UserRepository
}

func NewAdminService(ur repository.UserRepository) AdminService {
	return &adminService{
		UserRepository: ur,
	}
}

// FetchOrganizer implements AdminService.
func (s *adminService) FetchOrganizer() (model.ServiceResponse, error) {
	users, err := s.UserRepository.FindAllOrganizer()
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find organizer",
		}, err
	}

	if len(users) == 0 {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Record not found",
		}, err
	}

	var res []model.ProfileUserResponse

	for _, user := range users {
		dump := model.ProfileUserResponse{
			ID:                user.ID,
			Email:             user.Email,
			FullName:          user.FullName,
			NimNik:            user.NimNik,
			Prodi:             user.Prodi,
			Universitas:       user.Universitas,
			ID_Url:            user.IDUrl,
			IsOrganizer:       user.IsOrganizer,
			IsEmailVerified:   user.IsEmailVerified,
			IsProfileVerified: user.IsProfileVerified,
		}

		res = append(res, dump)
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find all organizer",
		Data:    res,
	}, nil
}

// FetchOrganizerDetail implements AdminService.
func (s *adminService) FetchOrganizerDetail(id string) (model.ServiceResponse, error) {
	user, err := s.UserRepository.FindById(id)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find organizer",
		}, err
	}

	res := model.ProfileUserResponse{
		ID:                user.ID,
		Email:             user.Email,
		FullName:          user.FullName,
		NimNik:            user.NimNik,
		Prodi:             user.Prodi,
		Universitas:       user.Universitas,
		ID_Url:            user.IDUrl,
		IsOrganizer:       user.IsOrganizer,
		IsEmailVerified:   user.IsEmailVerified,
		IsProfileVerified: user.IsProfileVerified,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully find organizer",
		Data:    res,
	}, nil
}

// FetchUser implements AdminService.
func (*adminService) FetchUser() (model.ServiceResponse, error) {
	panic("unimplemented")
}

// VerifyOrganizer implements AdminService.
func (s *adminService) VerifyOrganizer(req model.OrganizerVerifyRequest) (model.ServiceResponse, error) {
	user, err := s.UserRepository.FindById(req.ID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to find organizer",
		}, err
	}

	user.IsProfileVerified = req.Verify

	if req.Verify {
		// TODO Schedule to send email notification that user successfully created
		if err := s.UserRepository.Update(user); err != nil {
			return model.ServiceResponse{
				Code:    http.StatusInternalServerError,
				Error:   true,
				Message: "Something went wrong, failed to update user",
			}, err
		}
	} else {
		// TODO Schedule to send email notification that user is not eligible
		// TODO Schedule delete user
	}

	res := model.ProfileUserResponse{
		ID:                user.ID,
		Email:             user.Email,
		FullName:          user.FullName,
		NimNik:            user.NimNik,
		Prodi:             user.Prodi,
		Universitas:       user.Universitas,
		ID_Url:            user.IDUrl,
		IsOrganizer:       user.IsOrganizer,
		IsEmailVerified:   user.IsEmailVerified,
		IsProfileVerified: user.IsProfileVerified,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully verify organizer profile",
		Data:    res,
	}, nil
}
