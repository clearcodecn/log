package xlogger

import (
	"context"
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	AddHook(func(ctx context.Context) Field {
		reqid, ok := ctx.Value("reqid").(string)
		if !ok {
			return Field{}
		}
		return Any("reqid", reqid)
	})

	ctx := context.WithValue(context.Background(), "reqid", "123456")
	New(ctx).Info("help me")

	ctx = context.WithValue(context.Background(), "reqid", "1234")
	New(ctx).WithFields(Err(errors.New("some error"))).Error("help error")
}
