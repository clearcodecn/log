package xlogger

import (
	"context"
	"github.com/sirupsen/logrus"
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

func NewLog(config Config) (*logrus.Logger, error) {
	return newStdLog(config)
}

func SetGlobal(logger *zap.Logger) {
	std = logger
}

func newStdLog(config *Config) (*logrus.Logger, error) {
	l := logrus.New()
	if config.Level != "" {
		lvl, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}
		l.SetLevel(lvl)
	}
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

func Logger(ctx context.Context) *contextLogger {
	return &contextLogger{
		ctx:   ctx,
		log:   std,
		field: make(Field),
	}
}

func (c *contextLogger) New(ctx context.Context) *contextLogger {
	return Logger(ctx)
}

func (c *contextLogger) Debug(msg string) {
	fields := runHook(c)
	c.log.Debug(msg, fields...)
}

func (c *contextLogger) Info(msg string) {
	fields := runHook(c)
	c.log.Info(msg, fields...)
}

func (c *contextLogger) Error(msg string) {
	fields := runHook(c)
	c.log.Error(msg, fields...)
}

func (c *contextLogger) Warn(msg string) {
	fields := runHook(c)
	c.log.Warn(msg, fields...)
}

func (c *contextLogger) WithField(key string, val any) *contextLogger {
	c.field[key] = val
	return c
}

func (c *contextLogger) WithFields(fields ...Field) *contextLogger {
	for _, field := range fields {
		for k, v := range field {
			c.field[k] = v
		}
	}
	return c
}

func Debug(msg string, fields ...Field) {
	Logger(context.TODO()).WithFields(fields...).log.Debug(msg)
}

func Info(msg string, fields ...Field) {
	Logger(context.TODO()).WithFields(fields...).log.Debug(msg)
}

func Error(msg string, fields ...Field) {
	Logger(context.TODO()).WithFields(fields...).log.Debug(msg)
}

func Warn(msg string, fields ...Field) {
	Logger(context.TODO()).WithFields(fields...).log.Debug(msg)
}
