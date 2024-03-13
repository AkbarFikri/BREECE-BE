package service

import (
	"net/http"
	"strings"
	"time"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/AkbarFikri/BREECE-BE/internal/app/entity"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/gocron"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/mailer"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type UserService interface {
	Register(req model.CreateUserRequest) (model.ServiceResponse, error)
	Login(req model.LoginUserRequest) (model.ServiceResponse, error)
	VerifyOTP(req model.OtpUserRequest) (model.ServiceResponse, error)
	VerifyProfile(req model.ProfileUserRequest) (model.ServiceResponse, error)
	UserCurrent(req model.UserTokenData) (model.ServiceResponse, error)
}

type userService struct {
	UserRepository  repository.UserRepository
	CacheRepository repository.CacheRepository
	SupabaseBucket  *supabasestorageuploader.Client
}

func NewUserService(ur repository.UserRepository,
	cr repository.CacheRepository,
	sb *supabasestorageuploader.Client) UserService {
	return &userService{
		UserRepository:  ur,
		CacheRepository: cr,
		SupabaseBucket:  sb,
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

	user := entity.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: hashPass,
		FullName: req.FullName,
	}

	if err := s.UserRepository.Insert(user); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to create user",
			Data:    err.Error(),
		}, err
	}

	// Added OTP
	referenceID := helper.GenerateRandomString(16)
	OTP := helper.GenerateRandomInt(4)

	go mailer.Send(user.Email, "Your OTP Verification", string(OTP), user.FullName)

	s.CacheRepository.Set("otp:"+referenceID, []byte(OTP))
	s.CacheRepository.Set("user-ref:"+referenceID, []byte(user.ID))

	go gocron.ScheduleDeleteInvalidOtp(5*time.Minute, s.CacheRepository, referenceID)

	res := model.CreateUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		ReferenceID: referenceID,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully register user and send OTP to email",
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

	accessData := map[string]interface{}{
		"id":                user.ID,
		"email":             user.Email,
		"isAdmin":           user.IsAdmin,
		"isEmailVerified":   user.IsEmailVerified,
		"isProfileVerified": user.IsProfileVerified,
		"isOrganizer":       user.IsOrganizer,
		"isBrawijaya":       user.IsBrawijaya,
	}
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
			"token":     accessToken,
			"expire_at": time.Now().Add(24 * time.Hour).UTC().String(),
		},
	}, nil
}

func (s *userService) VerifyOTP(req model.OtpUserRequest) (model.ServiceResponse, error) {
	val, err := s.CacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid OTP",
			Data:    nil,
		}, err
	}

	otp := string(val)
	if otp != req.Otp {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid OTP",
			Data:    nil,
		}, err
	}

	val, _ = s.CacheRepository.Get("user-ref:" + req.ReferenceID)
	user, err := s.UserRepository.FindById(string(val))
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusNotFound,
			Error:   true,
			Message: "Something went wrong, unable to find user",
			Data:    nil,
		}, err
	}

	user.IsEmailVerified = true
	user.EmailVerifiedAt = time.Now()

	if err := s.UserRepository.Update(user); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, unable to update user",
			Data:    nil,
		}, err
	}

	data := model.OtpUserResponse{
		ID:              user.ID,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Successfully to verified email",
		Data:    data,
	}, nil
}

func (s *userService) VerifyProfile(req model.ProfileUserRequest) (model.ServiceResponse, error) {
	user, err := s.UserRepository.FindById(req.ID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "User with ID provided is not found.",
			Data:    nil,
		}, err
	}

	if user.ID_Url != "" {
		s.SupabaseBucket.Delete(user.ID_Url)
	}

	req.IdFile.Filename = "id_" + user.ID + ".png"

	idUrl, err := s.SupabaseBucket.Upload(req.IdFile)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to upload id_image to bucket",
			Data:    req.IdFile.Filename,
		}, err
	}

	user.ID_Url = idUrl
	user.NimNik = req.NimNik
	user.Prodi = req.Prodi
	user.Universitas = req.Universitas
	user.IsProfileVerified = true

	if strings.Contains(strings.ToLower(req.Universitas), "brawijaya") {
		user.IsBrawijaya = true
	}

	if err := s.UserRepository.Update(user); err != nil {
		return model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong, failed to update user profile",
			Data:    nil,
		}, err
	}

	res := model.ProfileUserResponse{
		ID:                user.ID,
		Email:             user.Email,
		FullName:          user.FullName,
		NimNik:            user.NimNik,
		Prodi:             user.Prodi,
		Universitas:       user.Universitas,
		ID_Url:            user.ID_Url,
		IsOrganizer:       user.IsOrganizer,
		IsEmailVerified:   user.IsEmailVerified,
		IsProfileVerified: user.IsProfileVerified,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Succesfully verify user profile.",
		Data:    res,
	}, nil
}

func (s *userService) UserCurrent(req model.UserTokenData) (model.ServiceResponse, error) {
	user, err := s.UserRepository.FindById(req.ID)
	if err != nil {
		return model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "User with ID provided is not found.",
			Data:    nil,
		}, err
	}

	res := model.ProfileUserResponse{
		ID:                user.ID,
		Email:             user.Email,
		FullName:          user.FullName,
		NimNik:            user.NimNik,
		Prodi:             user.Prodi,
		Universitas:       user.Universitas,
		ID_Url:            user.ID_Url,
		IsOrganizer:       user.IsOrganizer,
		IsEmailVerified:   user.IsEmailVerified,
		IsProfileVerified: user.IsProfileVerified,
	}

	return model.ServiceResponse{
		Code:    http.StatusOK,
		Error:   false,
		Message: "Succesfully find user profile.",
		Data:    res,
	}, nil
}
