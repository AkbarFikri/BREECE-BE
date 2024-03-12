package model

import (
	"mime/multipart"
	"time"
)

type FilterParam struct {
	Search   string
	Sort     string
	Page     int `default:1`
	Place    string
	Date     time.Time
	Category string
	IsPublic bool
}

type EventRequest struct {
	CategoryID   string                `form:"category_id" binding:"required"`
	Title        string                `form:"title" `
	Description  string                `form:"description"`
	Place        string                `form:"place" default:"Online"`
	Speakers     []string              `form:"speakers" binding:"required"`
	SpeakersRole []string              `form:"speaker_roles" binding:"required"`
	Banner       *multipart.FileHeader `form:"banner"`
	Date         string                `form:"date" binding:"required"`
	StartAt      string                `form:"start_at" binding:"required"`
	Link         string                `form:"link"`
	Price        uint32                `form:"price" default:0`
	TicketQty    uint16                `form:"ticket_qty" binding:"required"`
	IsPublic     bool                  `form:"is_public" binding:"required"`
}

type EventResponse struct {
	ID           string   `json:"-"`
	CategoryID   string   `json:"category_id"`
	Title        string   `json:"title" `
	Description  string   `json:"description"`
	Place        string   `json:"place"`
	Speakers     []string `json:"speakers"`
	SpeakersRole []string `json:"speaker_roles"`
	BannerUrl    string   `json:"banner"`
	Date         string   `json:"date"`
	StartAt      string   `json:"start_at"`
	Link         string   `json:"link"`
	Price        uint32   `json:"price"`
	TicketQty    uint16   `json:"ticket_qty"`
	OrganizeBy   string   `json:"organize_by"`
	IsPublic     bool     `json:"is_public"`
}
