package slogging

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

var TraceExemplar = func(ctx context.Context) prometheus.Labels {
	traceId, ok := ctx.Value(xb3traceid).(string)
	if !ok || traceId == "" {
		return nil
	}

	return prometheus.Labels{
		xb3traceid: traceId,
	}
}
