package slogging

import (
	"github.com/google/uuid"
)

func generateTraceId() string {
	traceID := uuid.New()
	return traceID.String()
}
