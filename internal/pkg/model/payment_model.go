package model

type PaymentRequest struct {
	EventID string `json:"event_id" binding:"required"`
}

type PaymentResponse struct {
	SnapUrl string `json:"snap_url"`
	Token   string `json:"token"`
}
