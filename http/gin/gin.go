package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scbt-ecom/slogging"
	"strconv"
)

func TraceMiddleware(l *slogging.SLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(slogging.XB3TraceID)
		if traceID == "" {
			traceID = slogging.GenerateTraceID()
		}

		var XB3Order int
		order := c.GetHeader(slogging.XB3Order)
		if order != "" {
			XB3Order, _ = strconv.Atoi(order)
		} else {
			XB3Order = -1
		}

		newL := l.Logger.With(slogging.StringAttr(slogging.XB3TraceID, traceID))
		sl := &slogging.SLogger{Logger: newL}

		ctx := slogging.ContextWithLogger(c.Request.Context(), sl, XB3Order)
		ctx = context.WithValue(ctx, slogging.XB3TraceID, traceID)

		c.Writer.Header().Set(slogging.XB3TraceID, traceID)
		c.Writer.Header().Set(slogging.XB3Order, strconv.Itoa(sl.GetOrder()))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
