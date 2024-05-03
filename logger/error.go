package logger

import (
	"log"
	"os"
)

var logError = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

func Error(message string, params interface{}) {
	logError.Println(message, params)
}
