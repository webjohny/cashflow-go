package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

func CreateHashByJson(state interface{}) string {
	output := JsonSerialize(state)
	hash := sha256.New()
	hash.Write([]byte(output)) // Assuming addSalt is a function that adds a salt
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
