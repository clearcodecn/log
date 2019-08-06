package log

import (
	"go.uber.org/zap"
	"testing"
)

func TestInfo(t *testing.T) {
	Info("foo", zap.String("key", "val"))
}

func TestFile(t *testing.T) {
	l, err := newStdLog(Config{
		File: "/tmp/im.log",
	})

	if err != nil {
		t.Fatal(err)
	}

	ReplaceLog(l)

	Info("foo", zap.String("file", "value"))
}
