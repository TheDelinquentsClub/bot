package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

//goland:noinspection GoUnusedGlobalVariable
type logger struct {
	Module string
	lock   sync.RWMutex
}

func init() {
	var logger = NewLogger("Logger")

	logger.Info("Test!")

}

func NewLogger(name string) *logger {
	return &logger{
		Module: name,
	}
}

func (l *logger) log(lvl string, args ...interface{}) {
	l.lock.RLock()

	t := time.Now().Format("02-01-2006")

	logFile, err := os.OpenFile(fmt.Sprintf("./logs/logs_%v.txt", t), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mv := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(mv)
	log.Println(fmt.Sprintf("%v %v: %v", l.Module, lvl, args))
	l.lock.RUnlock()
}

func (l *logger) Info(args ...interface{}) {
	l.log("Info", args)
}

func (l *logger) Warn(args ...interface{}) {
	l.log("Warn", args)
}

func (l *logger) Error(args ...interface{}) {
	l.log("Error", args)
}

func (l *logger) Debug(args ...interface{}) {
	l.log("Debug", args)
}

func (l *logger) Critical(args ...interface{}) {
	l.log("Critical", args)
	os.Exit(1)
}

func (l *logger) Notice(args ...interface{}) {
	l.log("Notice", args)
}
