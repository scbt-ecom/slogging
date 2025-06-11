package grpc

import (
	"context"
	"github.com/scbt-ecom/slogging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceInterceptor(l *slogging.SLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			traceID := slogging.GenerateTraceID()
			ctx = slogging.ContextWithLogger(ctx, &slogging.SLogger{l.With()})
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			traceId := GenerateTraceID()
			ctx = ContextWithLogger(ctx, &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
			ctx = context.WithValue(ctx, XB3TraceID, traceId)
			return handler(ctx, req)
		}

		traceIds := md.Get(XB3TraceID)
		var traceId string
		if len(traceIds) > 0 {
			traceId = traceIds[0]
		} else {
			traceId = GenerateTraceID()
		}

		ctx = ContextWithLogger(ctx, &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
		ctx = context.WithValue(ctx, XB3TraceID, traceId)
		return handler(ctx, req)
	}
}
