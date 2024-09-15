package gcp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// https://cloud.google.com/logging/docs/agent/configuration#special-fields
const (
	spanIDKey = "logging.googleapis.com/spanId"
	traceKey  = "logging.googleapis.com/trace"
)

func GetSpanField(ctx context.Context) zap.Field {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsSampled() {
		return zap.Skip()
	}
	return zap.String(spanIDKey, sc.SpanID().String())
}

func GetTraceField(ctx context.Context, projectID string) zap.Field {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsSampled() {
		return zap.Skip()
	}
	return zap.String(traceKey, fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID()))
}
