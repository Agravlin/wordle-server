package util

import (
	"math/rand"
)

func GenerateRoomCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	code := make([]byte, length)

	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}
