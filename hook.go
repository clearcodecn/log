package xlogger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Field map[string]any

type Hook func(ctx context.Context) Field

var globalHooks []Hook

// AddHook unsafe add global hook, it should be called at init functions.
func AddHook(hook Hook) {
	globalHooks = append(globalHooks, hook)
}

func runHook(c *contextLogger) logrus.Fields {
	for _, h := range globalHooks {
		field := h(c.ctx)
		for k, v := range field {
			c.field[k] = v
		}
	}
	var fields = make(logrus.Fields)
	for k, v := range c.field {
		fields[k] = v
	}
	return fields
}

func Err(err error) Field {
	return Field{
		"error": err,
	}
}

func Any(key string, val any) Field {
	return Field{
		key: val,
	}
}

// utils start

// reqidContext save a request id in context.
type reqidContext struct{}

func WithRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, reqidContext{}, id)
}

// ReqId set context with request id in all log trace.
func ReqId(ctx context.Context) Field {
	reqid, ok := ctx.Value(reqidContext{}).(string)
	if !ok {
		return Field{}
	}
	return Field{
		"reqid": reqid,
	}
}
