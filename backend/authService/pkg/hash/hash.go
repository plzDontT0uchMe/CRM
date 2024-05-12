package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateRandomHash() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("error generating random bytes: ", err)
	}

	hash := sha256.New()
	hash.Write(randomBytes)
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
