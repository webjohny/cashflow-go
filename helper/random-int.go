package helper

import (
	"math/rand"
	"time"
)

func Random(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max + 1)
}
