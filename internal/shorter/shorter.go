package shorter

import (
	"crypto/rand"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShort() string {
	b := make([]byte, 10)
	rand.Read(b)
	result := make([]byte, 10)
	for i := range b {
		result[i] = symbols[b[i]%byte(len(symbols))]
	}
	return string(result)
}
