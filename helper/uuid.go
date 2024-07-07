package helper

import "github.com/google/uuid"

func Uuid(str string) string {
	return str + "_" + uuid.New().String()
}
