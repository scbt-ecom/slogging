package slogging

import (
	"context"
	"net/http"
)

type HTTPMiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

//func HTTPTraceMiddleware(logger *slog.Logger) HTTPMiddlewareFunc {
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			traceId := r.Header.Get(XB3TraceID)
//			if traceId == "" {
//				traceId = GenerateTraceID()
//			}
//
//			ctx := ContextWithLogger(r.Context(), &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
//			ctx = context.WithValue(ctx, XB3TraceID, traceId)
//
//			next.ServeHTTP(w, r.WithContext(ctx))
//		}
//	}
//}
//
//func MuxHTTPTraceMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			traceId := r.Header.Get(XB3TraceID)
//			if traceId == "" {
//				traceId = GenerateTraceID()
//			}
//
//			ctx := ContextWithLogger(r.Context(), &SLogger{logger.With(StringAttr(XB3TraceID, traceId))})
//			ctx = context.WithValue(ctx, XB3TraceID, traceId)
//
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
//}

func RequestWithTraceHeaders(ctx context.Context, req *http.Request) *http.Request {
	traceId := ctx.Value(XB3TraceID).(string)

	req.Header.Set(XB3TraceID, traceId)
	return req
}
