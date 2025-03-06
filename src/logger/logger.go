package logger

import (
	"first-project/src/bootstrap"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Logger struct {
	zap *zap.Logger
}

var loggerInstance *Logger

func NewLogger(config bootstrap.LoggerConfig) (*Logger, error) {
	if config.LogLevel == "" {
		config.LogLevel = InfoLevel
	}

	var level zapcore.Level
	switch config.LogLevel {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleOutput, err := strconv.ParseBool(config.ConsoleOutput)
	if err != nil {
		consoleOutput = true
		log.Println("Error during checking console output enable. set it to true by default")
	}
	if consoleOutput {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(file),
			level,
		)
		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger := &Logger{
		zap: zapLogger,
	}

	loggerInstance = logger
	return logger, nil
}

func GetLogger() *Logger {
	if loggerInstance == nil {
		loggerInstance = &Logger{
			zap: zap.NewExample(),
		}
	}
	return loggerInstance
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &Logger{
		zap: l.zap.With(zapFields...),
	}
}

func (l *Logger) Close() {
	_ = l.zap.Sync()
}
