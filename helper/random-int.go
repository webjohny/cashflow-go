package helper

import (
	"math/rand"
	"time"
)

func Random(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Инициализация генератора случайных чисел
	return r.Intn(max + 1)
}

func RandomMinMax(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
