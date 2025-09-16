package random

import (
	"math/rand"
	"time"
)

func NewRandomString(prefix string, size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	var res = make([]rune, size)

	for i := range res {
		res[i] = chars[rnd.Intn(len(chars))]
	}

	return prefix + string(res)
}