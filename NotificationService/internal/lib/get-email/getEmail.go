package getemail

import (
	"encoding/json"

	"github.com/LashkaPashka/notification-service/internal/model"
)


func GetEmail(payload []byte) string{
	var data model.OrderPaid

	err := json.Unmarshal(payload, &data)
	if err != nil {
		panic(err)
	}
	
	return data.DataM.Email
}