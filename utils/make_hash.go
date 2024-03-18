package utils

import (
	"crypto/sha256"
	"fmt"
)

func MakeHash(str string) (hashedStr string) {
	hash := sha256.Sum256([]byte(str))
	hashedStr = fmt.Sprintf("%x", hash)

	return hashedStr
}
