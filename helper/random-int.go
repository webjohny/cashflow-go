package helper

import (
	"math/rand"
	"time"
)

func Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Инициализация генератора случайных чисел
	return r.Intn(max + 1)
}
