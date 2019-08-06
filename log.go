package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	std *zap.Logger
)

func init() {
	var err error
	std, err = newStdLog(Config{})
	if err != nil {
		panic(err)
	}
}

const (
	defaultLogSize = 300 * (2 << 19)
)

type Config struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	File   string `json:"file"`
}

func NewLog(config Config) (*zap.Logger, error) {
	return newStdLog(config)
}

func newStdLog(config Config) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()

	enc := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "t",
		LevelKey:       "l",
		NameKey:        "n",
		CallerKey:      "c",
		MessageKey:     "m",
		StacktraceKey:  "s",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncode(),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if config.File != "" {
		zapConfig.OutputPaths = []string{config.File}
	} else {
		zapConfig.Encoding = "console"

	}
	zapConfig.EncoderConfig = enc

	var level = zapcore.InfoLevel
	if config.Level != "" {
		_ = level.Set(config.Level)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	return zapConfig.Build(zap.AddCallerSkip(1),
		zap.AddCaller())
}

func timeEncode() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("01/02 15:04:05"))
	}
}

func Debug(msg string, fields ...zap.Field) {
	std.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	std.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	std.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	std.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	std.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	std.Panic(msg, fields...)
}

func ReplaceLog(l *zap.Logger) {
	_ = std.Sync()
	*std = *l
}
