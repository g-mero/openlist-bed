package logger

import "go.uber.org/zap/zapcore"

type Config struct {
	Level          string
	FilePath       string
	DisableConsole bool
	DbCore         zapcore.Core
}
