package payload

import "time"


type PostEventRequest struct {
	Title string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Tags []string `json:"tags"`
	Timestamp time.Time
}

type LikeEventRequest struct {
	Timestamp time.Time
}

type OrderEventRequest struct {
	Amount int `json:"amount" validate:"required"`
	Currency string `json:"currency" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	Timestamp time.Time
}