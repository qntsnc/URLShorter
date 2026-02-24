package shorter

import (
	"crypto/sha256"
	"encoding/binary"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShort(url string) string {
	hash := sha256.Sum256([]byte(url))
	num := binary.BigEndian.Uint64(hash[:8])
	base := uint64(len(symbols))
	var result []byte
	for i := 0; i < 10; i++ {
		rem := num % base
		result = append(result, symbols[rem])
		num = num / base
	}
	return string(result)
}
