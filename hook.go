package xlogger

import (
	"context"
	"go.uber.org/zap"
)

type Field map[string]any

type Hook func(ctx context.Context) Field

var globalHooks []Hook

// AddHook unsafe add global hook, it should be called at init functions.
func AddHook(hook Hook) {
	globalHooks = append(globalHooks, hook)
}

func runHook(ctx context.Context) []zap.Field {
	var fields = make(Field)
	for _, h := range globalHooks {
		field := h(ctx)
		for k, v := range field {
			fields[k] = v
		}
	}
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
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
