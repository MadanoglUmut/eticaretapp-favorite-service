package logs

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
}

func Info(message string, opts ...func(*slog.Logger) *slog.Logger) {
	l := logger
	for _, opt := range opts {
		l = opt(l)
	}
	l.Info(message)
}

func Error(message string, opts ...func(*slog.Logger) *slog.Logger) {
	l := logger
	for _, opt := range opts {
		l = opt(l)
	}
	l.Error(message)
}

func Warning(message string, opts ...func(*slog.Logger) *slog.Logger) {
	l := logger
	for _, opt := range opts {
		l = opt(l)
	}
	l.Warn(message)
}

func WithHandlerName(name string) func(*slog.Logger) *slog.Logger {
	return func(l *slog.Logger) *slog.Logger {
		return l.With("Handler", name)
	}
}

func WithStatus(status int) func(*slog.Logger) *slog.Logger {
	return func(l *slog.Logger) *slog.Logger {
		return l.With("Status", status)
	}
}
