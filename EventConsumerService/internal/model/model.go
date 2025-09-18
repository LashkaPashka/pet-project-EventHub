package model

import (
	"time"

	modelforevents "github.com/LashkaPashka/EventConsumerService/internal/model/model-for-events"
)

type eventMetaData struct {
	ID string 									`json:"id"`
	Type string 								`json:"type"`
	Timestamp time.Time 						`json:"timestamp"`
	Source string 								`json:"source"`
	MetaM Meta 									`json:"meta"`
}

type UserPostCreated struct {
	eventMetaData
	DataM modelforevents.PostData 				`json:"data"`
}

type UserPostLiked struct {
	eventMetaData
	DataM modelforevents.LikeData 				`json:"data"`
}

type OrderPaid struct {
	eventMetaData
	DataM modelforevents.OrderData 				`json:"data"`
}


type Meta struct {
	Trace_id string	`json:"trace_id"`
	Version string  `json:"version"`
}