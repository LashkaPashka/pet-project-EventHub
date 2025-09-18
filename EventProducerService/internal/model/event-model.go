package model

import (
	"time"
)

type EventForKafka struct {
	ID string 						`json:"id"`
	Type string 					`json:"type"`
	Timestamp time.Time 			`json:"timestamp"`
	Source string 					`json:"source"`
	DataM any 						`json:"data"`
	MetaM Meta 						`json:"meta"`
}

type Meta struct {
	Trace_id string `json:"trace_id"`
	Version string  `json:"version"`
}