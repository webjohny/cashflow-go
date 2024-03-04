package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

func CreateHash(state string) string {
	hash := sha256.New()
	hash.Write([]byte(state)) // Assuming addSalt is a function that adds a salt
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

func CreateHashByJson(state interface{}) string {
	output := JsonSerialize(state)
	return CreateHash(output)
}
