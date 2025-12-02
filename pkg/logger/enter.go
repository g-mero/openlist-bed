package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// Init 初始化日志
func Init(cfg Config) error {
	Logger = NewZapLogger(cfg)
	return nil
}

func selectEncoder(encoding string, cfg zapcore.EncoderConfig) zapcore.Encoder {
	if strings.ToLower(encoding) == "json" {
		return zapcore.NewJSONEncoder(cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}

func Debug(msg string, fields ...zap.Field) { Logger.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { Logger.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Logger.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Logger.Error(msg, fields...) }

func NewZapLogger(cfg Config) *zap.Logger {
	var level zapcore.Level
	switch strings.ToUpper(cfg.Level) {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "msg"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var cores []zapcore.Core

	if !cfg.DisableConsole {
		cores = append(cores, zapcore.NewCore(selectEncoder("console", encoderCfg), zapcore.AddSync(os.Stdout), level))
	}

	if cfg.FilePath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    25, // MB
			MaxBackups: 6,
			MaxAge:     30, // days
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(selectEncoder("console", encoderCfg), fileWriter, level))
	}

	if cfg.DbCore != nil {
		cores = append(cores, cfg.DbCore)
	}

	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(1)}
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...)
}
