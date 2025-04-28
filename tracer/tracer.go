package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const traceKey contextKey = "traceRecorder"

type TraceRecorder struct {
	TraceID   string
	StartTime time.Time
	Events    []string
}

func (t *TraceRecorder) Record(entry string) {
	if t == nil {

		t = NewTraceRecorder(uuid.New().String())
	}
	t.Events = append(t.Events,
		fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339),
			entry))
	fmt.Println("printing ", t.Events[len(t.Events)-1])
}

func NewTraceRecorder(traceID string) *TraceRecorder {
	return &TraceRecorder{
		TraceID:   traceID,
		StartTime: time.Now(),
		Events:    []string{},
	}
}

func (r *TraceRecorder) GetTraceID() string {
	return r.TraceID
}

func FromContext(ctx context.Context) *TraceRecorder {
	traceRecorder, _ := ctx.Value(traceKey).(*TraceRecorder)
	if traceRecorder == nil {
		traceID := uuid.New().String()
		return NewTraceRecorder(traceID)
	}
	return traceRecorder
}

func ToContext(ctx context.Context, traceRecorder *TraceRecorder) context.Context {
	return context.WithValue(ctx, traceKey, traceRecorder)
}
