package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type UserService interface {
	Register(req model.CreateUserRequest) (model.ServiceResponse, error)
	Login(req model.LoginUserRequest) (model.ServiceResponse, error)
}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		UserRepository: ur,
	}
}

func (s *userService) Register(req model.CreateUserRequest) (model.ServiceResponse, error) {
	_, err := s.UserRepository.FindByEmail(req.Email)
	if err == nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Email already used",
			Data:    nil,
		}, err
	}

	hashPass, err := helper.HashPassword(req.Password)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong",
			Data:    err.Error(),
		}, err
	}

	var user entity.User
	user.ID = uuid.New().String()
	user.Email = req.Email
	user.Username = req.Username
	user.Password = hashPass

	if err := s.UserRepository.Insert(user); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to create user",
			Data:    err.Error(),
		}, err
	}

	res := model.CreateUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully register user",
		Data:    res,
	}, nil
}

func (s *userService) Login(req model.LoginUserRequest) (model.ServiceResponse, error) {
	user, err := s.UserRepository.FindByEmail(req.Email)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid Email or Password",
			Data:    err.Error(),
		}, err
	}

	if err := helper.ComparePassword(user.Password, req.Password); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid Email or Password",
			Data:    nil,
		}, err
	}

	accessData := map[string]interface{}{"id": user.ID, "email": user.Email}
	accessToken, err := helper.SignJWT(accessData, 24)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong",
			Data:    nil,
		}, err
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully login",
		Data: gin.H{
			"token":    accessToken,
			"expireAt": time.Now().Add(24 * time.Hour),
		},
	}, nil
}
