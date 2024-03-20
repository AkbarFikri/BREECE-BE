package model

import (
	"mime/multipart"
	"time"
)

type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required"`
	IsOrganizer bool   `json:"-"`
	// Nim         uint64 `json:"nim" binding:"required"`
	// Name        string `json:"name" binding:"required"`
	// Prodi       string `json:"prodi" binding:"required"`
	// Universitas string `json:"universitas" binding:"required"`
}

type CreateUserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	ReferenceID string `json:"reference_id"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type OtpUserRequest struct {
	ReferenceID string `json:"reference_id" binding:"required"`
	Otp         string `json:"otp" binding:"required"`
}

type OtpUserResponse struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
}

type ProfileUserRequest struct {
	ID          string                `form:"id" binding:"required"`
	NimNik      string                `form:"nim_nik" binding:"required"`
	Prodi       string                `form:"prodi"`
	Universitas string                `form:"universitas" binding:"required"`
	IdFile      *multipart.FileHeader `form:"id_file" binding:"required"`
}

type ProfileUserResponse struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	FullName          string `json:"full_name"`
	NimNik            string `json:"nim_nik"`
	Prodi             string `json:"prodi"`
	Universitas       string `json:"universitas"`
	ID_Url            string `json:"id_url"`
	IsOrganizer       bool   `json:"is_organizer"`
	IsEmailVerified   bool   `json:"is_email_verified"`
	IsProfileVerified bool   `json:"is_profile_verified"`
}

type UserTokenData struct {
	ID                string
	Email             string
	IsEmailVerified   bool
	IsProfileVerified bool
	IsAdmin           bool
	IsOrganizer       bool
	IsBrawijaya       bool
}

type OrganizerVerifyRequest struct {
	ID     string `json:"id"`
	Verify bool   `json:"verify"`
}
