package model

type EmailOTP struct {
	Subject string
	Email   string
	Otp     string
	Name    string
}

type EmailNotification struct {
	Subject    string
	Email      string
	Name       string
	EventTitle string
	EventStart string
	Venue      string
}

type EmailApproval struct {
	Subject string
	Email   string
	Name    string
	Status  string
}
