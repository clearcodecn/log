package xlogger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"strings"
	"time"
)

var (
	std        *logrus.Logger
	callPrefix string
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
	Level      string `json:"level"`
	Format     string `json:"format"`
	File       string `json:"file"`
	CallPrefix string `json:"callerPrefix"`
}

func NewLog(config Config) (*logrus.Logger, error) {
	return newStdLog(config)
}

func SetGlobal(logger *logrus.Logger) {
	std = logger
}

func newStdLog(config Config) (*logrus.Logger, error) {
	l := logrus.New()
	if config.Level != "" {
		lvl, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}
		l.SetLevel(lvl)
	}
	if config.File != "" {
		l.SetOutput(&lumberjack.Logger{
			Filename:   config.File,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		})
	}

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   `060102-150405`,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	callPrefix = config.CallPrefix
	return l, nil
}

func timeEncode() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
}

type contextLogger struct {
	ctx   context.Context
	field Field
	log   *logrus.Logger
}

func WithContext(ctx context.Context) *contextLogger {
	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		ctx = ginCtx.Request.Context()
	}
	return &contextLogger{
		ctx:   ctx,
		log:   std,
		field: make(Field),
	}
}

func (c *contextLogger) New(ctx context.Context) *contextLogger {
	newLogger := WithContext(ctx)
	return newLogger.WithFields(c.field)
}

func (c *contextLogger) Debug(msg string) {
	fields := runHook(c)
	caller, line := c.getCaller(1)
	if caller != "" {
		fields["caller"] = fmt.Sprintf("%s:%d", caller, line)
	}
	c.log.WithFields(fields).Debug(msg)
}

func (c *contextLogger) Info(msg string) {
	fields := runHook(c)
	caller, line := c.getCaller(1)
	if caller != "" {
		fields["caller"] = fmt.Sprintf("%s:%d", caller, line)
	}
	c.log.WithFields(fields).Info(msg)
}

func (c *contextLogger) Error(msg string) {
	fields := runHook(c)
	caller, line := c.getCaller(1)
	if caller != "" {
		fields["caller"] = fmt.Sprintf("%s:%d", caller, line)
	}
	c.log.WithFields(fields).Error(msg)
}

func (c *contextLogger) Warn(msg string) {
	fields := runHook(c)
	caller, line := c.getCaller(1)
	if caller != "" {
		fields["caller"] = fmt.Sprintf("%s:%d", caller, line)
	}
	c.log.WithFields(fields).Warn(msg)
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

func Debug(ctx context.Context, msg string, fields ...Field) {
	WithContext(ctx).WithFields(fields...).log.Debug(msg)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	WithContext(ctx).WithFields(fields...).log.Info(msg)
}

func Error(ctx context.Context, msg string, fields ...Field) {
	WithContext(ctx).WithFields(fields...).log.Error(msg)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	WithContext(ctx).WithFields(fields...).log.Warn(msg)
}

func WithErr(ctx context.Context, err error) *contextLogger {
	return WithContext(ctx).WithFields(Err(err))
}

// logWithCaller 记录日志并包含调用者信息
func (l *contextLogger) getCaller(n int) (string, int) {
	_, file, line, _ := runtime.Caller(n + 1) // 获取调用者信息
	if len(file) != 0 {
		file = strings.TrimPrefix(file, callPrefix)
	}
	return file, line
}
