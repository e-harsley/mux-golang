package happiness

import (
	"fmt"
	"math/rand"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateSecretKey(test bool) string {
	prefix := "sk_live"
	if test {
		prefix = "sk_test"
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	randomBytes := make([]byte, 18) // 18 bytes for base64 encoding to result in 24 characters
	rand.Read(randomBytes)
	for i := range randomBytes {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return fmt.Sprintf("%s_%s", prefix, string(randomBytes))
}
