package slogging

import (
	"context"
	"log/slog"
)

type ctxLogger struct{}

func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func Context() context.Context {
	traceID := generateTraceId()

	l := slog.Default().With(StringAttr(xb3traceid, traceID))
	ctx := context.WithValue(context.Background(), xb3traceid, traceID)
	return context.WithValue(ctx, ctxLogger{}, l)
}

func L(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*Logger); ok {
		return l
	}

	if traceId, ok := ctx.Value(xb3traceid).(string); ok {
		return slog.Default().With(StringAttr(xb3traceid, traceId))
	}

	return slog.Default()
}
