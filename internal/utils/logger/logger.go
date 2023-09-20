package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

func (l *Level) String() string {
	return string(*l)
}

const (
	DebugLevel Level = "debug"
	// InfoLevel is the default logging priority.
	InfoLevel = "info"
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = "warn"
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = "error"
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = "dpanic"
	// PanicLevel logs a message, then panics.
	PanicLevel = "panic"
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = "fatal"
)

var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level Level) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level.String())
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}
