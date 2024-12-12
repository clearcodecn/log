package xlogger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	std *zap.Logger
)

func init() {
	var err error
	std, err = newStdLog(defaultConfig)
	if err != nil {
		panic(err)
	}
}

var (
	defaultConfig = Config{
		Level:  "info",
		Format: "json",
		File:   "",
	}
)

type Config struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	File   string `json:"file"`
}

func NewLog(config Config) (*zap.Logger, error) {
	return newStdLog(config)
}

func SetGlobal(logger *zap.Logger) {
	std = logger
}

func newStdLog(config Config) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()

	enc := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "ti",
		LevelKey:       "lvl",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
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

	var level = zapcore.DebugLevel
	if config.Level != "" {
		_ = level.Set(config.Level)
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	return zapConfig.Build(zap.AddCallerSkip(1), zap.AddCaller())
}

func timeEncode() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
}

type contextLogger struct {
	ctx   context.Context
	field Field
	log   *zap.Logger
}

func New(ctx context.Context) *contextLogger {
	return &contextLogger{
		ctx:   ctx,
		log:   std,
		field: make(Field),
	}
}

func (c *contextLogger) New(ctx context.Context) *contextLogger {
	return New(ctx)
}

func (c *contextLogger) Debug(msg string) {
	fields := runHook(c.ctx)
	c.log.Debug(msg, fields...)
}

func (c *contextLogger) Info(msg string) {
	fields := runHook(c.ctx)
	c.log.Info(msg, fields...)
}

func (c *contextLogger) Error(msg string) {
	fields := runHook(c.ctx)
	c.log.Error(msg, fields...)
}

func (c *contextLogger) Warn(msg string) {
	fields := runHook(c.ctx)
	c.log.Warn(msg, fields...)
}

func (c *contextLogger) WithField(key string, val any) *contextLogger {
	c.field[key] = val
	return c
}

func (c *contextLogger) WithFields(fields Field) *contextLogger {
	for k, v := range fields {
		c.field[k] = v
	}
	return c
}

func Debug(msg string, field Field) {
	New(context.TODO()).WithFields(field).log.Debug(msg)
}

func Info(msg string, field Field) {
	New(context.TODO()).WithFields(field).log.Debug(msg)
}

func Error(msg string, field Field) {
	New(context.TODO()).WithFields(field).log.Debug(msg)
}

func Warn(msg string, field Field) {
	New(context.TODO()).WithFields(field).log.Debug(msg)
}
