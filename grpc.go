package slogging

import (
	"context"
	"google.golang.org/grpc/metadata"
)

//func GRPCTraceMiddleware(logger *slog.Logger) grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//		md, ok := metadata.FromIncomingContext(ctx)
//		if !ok {
//			traceId := GenerateTraceID()
//			ctx = ContextWithLogger(ctx, &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
//			ctx = context.WithValue(ctx, XB3TraceID, traceId)
//			return handler(ctx, req)
//		}
//
//		traceIds := md.Get(XB3TraceID)
//		var traceId string
//		if len(traceIds) > 0 {
//			traceId = traceIds[0]
//		} else {
//			traceId = GenerateTraceID()
//		}
//
//		ctx = ContextWithLogger(ctx, &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
//		ctx = context.WithValue(ctx, XB3TraceID, traceId)
//		return handler(ctx, req)
//	}
//}

// example
//func GRPCExampleUsage() {
//	log := NewLogger(
//		SetLevel("debug"),
//		InGraylog("graylog:12201", "debug", "application_name"),
//		SetDefault(true),
//	)
//
//	tracemiddleware := GRPCTraceMiddleware(log)
//
//	srv := grpc.NewServer(
//		grpc.UnaryInterceptor(tracemiddleware),
//	)
//}

func MetadataWithTraceHeaders(ctx context.Context) context.Context {
	traceId, ok := ctx.Value(XB3TraceID).(string)
	if !ok || traceId == "" {
		traceId = GenerateTraceID()
	}

	md := metadata.Pairs(XB3TraceID, traceId)
	return metadata.NewOutgoingContext(context.Background(), md)
}

// example
//func exampleMetadataSend(ctx context.Context) {
//	mdctx := MetadataWithTraceHeaders(ctx)
//	client.HelloWorld(ctx, req)
//}
