package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {
	t := time.Now().Format("02-01-2006")

	logFile, err := os.OpenFile(fmt.Sprintf("./logs/logs_%v.txt", t), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Info = log.New(logFile, "[INFO]: ", log.Ltime|log.Lshortfile)
	Warn = log.New(logFile, "[WARN]: ", log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "[ERROR]: ", log.Ltime|log.Lshortfile)
}
