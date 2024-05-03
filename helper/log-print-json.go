package helper

import "log"

func LogPrintJson(state interface{}) {
	log.Println(JsonSerialize(state))
}
