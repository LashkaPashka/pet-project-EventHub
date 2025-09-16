package modelforevents

import "time"

type OrderData struct {
	Email string  `json:"email"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Timestamp time.Time `json:"timestamp"`
}