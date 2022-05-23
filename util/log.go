package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	consoleEncoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       zapcore.OmitKey,
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     zapcore.OmitKey,
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: zapcore.OmitKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
	})
	logOut    = zapcore.AddSync(os.Stdout)
	levelFunc = zap.DebugLevel /*zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})*/
	LoggerCore = zapcore.NewCore(consoleEncoder, logOut, levelFunc)
)
