package model

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	// Nim         uint64 `json:"nim" binding:"required"`
	// Name        string `json:"name" binding:"required"`
	// Prodi       string `json:"prodi" binding:"required"`
	// Universitas string `json:"universitas" binding:"required"`
}

type CreateUserResponse struct {
	ID          string
	Email       string
	FullName    string
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
	ID              string
	Email           string
	IsEmailVerified bool
	EmailVerifiedAt time.Time
}
