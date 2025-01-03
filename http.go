package slogging

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type HTTPMiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func HTTPTraceMiddleware(logger *slog.Logger) HTTPMiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			traceId := r.Header.Get(xb3traceid)
			if traceId == "" {
				traceId = generateTraceId()
			}

			ctx := ContextWithLogger(r.Context(), logger.With(StringAttr(xb3traceid, traceId)))
			ctx = context.WithValue(ctx, xb3traceid, traceId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func MuxHTTPTraceMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceId := r.Header.Get(xb3traceid)
			if traceId == "" {
				traceId = generateTraceId()
			}

			ctx := ContextWithLogger(r.Context(), logger.With(StringAttr(xb3traceid, traceId)))
			ctx = context.WithValue(ctx, xb3traceid, traceId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GinHTTPTraceMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(xb3traceid)
		if traceId == "" {
			traceId = generateTraceId()
		}

		ctx := ContextWithLogger(c.Request.Context(), logger.With(StringAttr(xb3traceid, traceId)))
		ctx = context.WithValue(ctx, xb3traceid, traceId)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestWithTraceHeaders(ctx context.Context, req *http.Request) *http.Request {
	traceId := ctx.Value(xb3traceid).(string)

	req.Header.Set(xb3traceid, traceId)
	return req
}
