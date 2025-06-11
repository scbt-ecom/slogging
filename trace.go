package slogging

import (
	"github.com/google/uuid"
)

func GenerateTraceID() string {
	traceID := uuid.New()
	return traceID.String()
}
