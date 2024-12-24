package util

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)

	for i := 0; i < length; i++ {
		randomIndex := randGen.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
