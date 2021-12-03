package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var LogFile *os.File

func GetLogWriter() zapcore.WriteSyncer {
	path, err := os.Getwd()
	if err != nil {
		Panic(err)
	}

	_, err = os.Stat("./logs")
	if os.IsNotExist(err) {
		_ = os.Mkdir("./logs", os.ModePerm)
	}

	t := time.Now().Format("02-01-2006")

	LogFile, err = os.OpenFile(fmt.Sprintf("%v/logs/logs_%v.txt", path, t), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	syncs := zapcore.NewMultiWriteSyncer(LogFile, os.Stdout)

	return syncs
}

func GetEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("02/01/2006 || 15:04:05"))
	}

	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(conf)
}

func Print(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	print(argsfmt)
}

func Debug(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Debug(argsfmt)
}

func Info(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Info(argsfmt)
}

func Warn(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Warn(argsfmt)
}

func Error(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Error(argsfmt)
}

func Panic(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Panic(argsfmt)
}

func Fatal(args ...interface{}) {
	argsfmt := strings.ReplaceAll(fmt.Sprint(args), "[", "")
	argsfmt = strings.ReplaceAll(argsfmt, "]", "")
	zap.S().Fatal(argsfmt)
}
