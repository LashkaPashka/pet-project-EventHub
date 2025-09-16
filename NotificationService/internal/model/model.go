package model

import (
	"time"

	modelforevents "github.com/LashkaPashka/notification-service/internal/model/model-for-events"
)

type OrderPaid struct {
	ID string 									`json:"id"`
	Type string 								`json:"type"`
	Timestamp time.Time 						`json:"timestamp"`
	Source string 								`json:"source"`
	DataM modelforevents.OrderData 				`json:"data"`
	MetaM Meta 									`json:"meta"`
}

type Meta struct {
	Trace_id string	`json:"trace_id"`
	Version string  `json:"version"`
}