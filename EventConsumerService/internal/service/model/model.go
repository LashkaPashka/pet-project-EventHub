package model

import (
	"time"

	modelforevents "github.com/LashkaPashka/EventConsumerService/internal/service/model/model-for-events"
)

type UserPostCreated struct {
	ID string 									`json:"id"`
	Type string 								`json:"type"`
	Timestamp time.Time 						`json:"timestamp"`
	Source string 								`json:"source"`
	DataM modelforevents.PostData 				`json:"data"`
	MetaM Meta 									`json:"meta"`
}

type UserPostLiked struct {
	ID string 									`json:"id"`
	Type string 								`json:"type"`
	Timestamp time.Time 						`json:"timestamp"`
	Source string 								`json:"source"`
	DataM modelforevents.LikeData 				`json:"data"`
	MetaM Meta 									`json:"meta"`
}

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