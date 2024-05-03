package logger

import (
	"log"
	"os"
)

var logInfo = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

func Info(message string, params interface{}) {
	logInfo.Println(message, params)
}
