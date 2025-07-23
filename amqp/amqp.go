package amqp

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"github.com/scbt-ecom/slogging"
	"github.com/scbt-ecom/rabbitmq"
	"log/slog"
)

func TraceMiddleware(l *slog.Logger) func(context.Context, amqp091.Delivery) context.Context {
	return func(baseCtx context.Context, msg amqp091.Delivery) context.Context {
		traceID, ok := msg.Headers[slogging.XB3TraceID].(string)
		if !ok || traceID == "" {
			traceID = slogging.GenerateTraceID()
		}

		newL := l.With(slogging.StringAttr(slogging.XB3TraceID, traceID))
		sl := &slogging.SLogger{Logger: newL}

		ctx := slogging.ContextWithLogger(baseCtx, sl)
		ctx = context.WithValue(ctx, slogging.XB3TraceID, traceID)

		return ctx
	}
}

func TraceHeaders(ctx context.Context, headers rabbitmq.Headers) rabbitmq.Headers {
	traceID, ok := ctx.Value(slogging.XB3TraceID).(string)
	if !ok || traceID == "" {
		traceID = slogging.GenerateTraceID()
	}

	headers[slogging.XB3TraceID] = traceID
	return headers
}
