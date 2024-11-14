package helper

import (
	"strconv"
)

func ConvertToBool(value string) bool {
	return value == "true"
}

func ConvertToInt(value string) int {
	valueParam, _ := strconv.Atoi(value)
	return valueParam
}

func ConvertToUInt8(value string) uint8 {
	valueParam, _ := strconv.Atoi(value)
	return uint8(valueParam)
}

func ConvertToUInt32(value string) uint32 {
	valueParam, _ := strconv.Atoi(value)
	return uint32(valueParam)
}

func ConvertToUInt64(value string) uint64 {
	valueParam, _ := strconv.Atoi(value)
	return uint64(valueParam)
}
