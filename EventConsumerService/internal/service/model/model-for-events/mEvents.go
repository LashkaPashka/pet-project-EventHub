package modelforevents

import "time"

type PostData struct {
	Email string  `json:"email"`
	Title string   `json:"title"`
	Content string `json:"content"`
	Tags []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

type LikeData struct {
	Email string  `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

type OrderData struct {
	Email string  `json:"email"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Timestamp time.Time `json:"timestamp"`
}