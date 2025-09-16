package encode

import "encoding/json"

func EncodeBytes(payload any) []byte {
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	return bytes
}