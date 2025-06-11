package slogging

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type AMQPMiddlewareFunc func(amqp091.Delivery) context.Context

//func AMQPTraceMiddleware(l *Logger) AMQPMiddlewareFunc {
//	return func(msg amqp091.Delivery) context.Context {
//		traceID, ok := msg.Headers[XB3TraceID].(string)
//		if !ok || traceID == "" {
//			traceID = GenerateTraceID()
//		}
//
//		ctx := ContextWithLogger(context.Background(), &SLogger{l.With(StringAttr(XB3TraceID, traceID))})
//		ctx = context.WithValue(ctx, XB3TraceID, traceID)
//
//		return ctx
//	}
//}

// example
//func ExampleAMQPTracing() {
//	log := NewLogger(
//		SetLevel("debug"),
//		InGraylog("graylog:12201", "debug", "application_name"),
//		SetDefault(true),
//	)
//
//	msgs := make(<-chan amqp091.Delivery)
//	amqpTraceMiddleware := AMQPTraceMiddleware(log)
//
//	go func() {
//		for msg := range msgs {
//			ProcessMessage(amqpTraceMiddleware(msg))
//		}
//	}()
//
//}
//
//func ProcessMessage(ctx context.Context, msg amqp091.Delivery) bool {
//	// SOME LOGIC
//	return true
//}

func AMQPTableWithTraceHeaders(ctx context.Context, table amqp091.Table) amqp091.Table {
	traceID, ok := ctx.Value(XB3TraceID).(string)
	if !ok || traceID == "" {
		traceID = GenerateTraceID()
	}

	table[XB3TraceID] = traceID
	return table
}
