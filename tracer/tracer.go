package tracer

import (
	"context"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/p4-pankaj/trace-replay/db"
	"github.com/p4-pankaj/trace-replay/models"
	"github.com/p4-pankaj/trace-replay/utility"
	"github.com/rs/zerolog"
)

type contextKey string

const traceKey contextKey = "traceRecorder"

type TraceRecorder struct {
	TraceID      string
	StartTime    time.Time
	traceRecord  *models.TraceRecord
	Updater      chan func(r *models.TraceRecord) //// TODO ,explore reflect package to decide how updater should work.
	CloseUpdater chan struct{}
	Lg           zerolog.Logger
	db           db.TraceStorage
}

// Warning: Potential For Goroutine Leak.
func (t *TraceRecorder) recorder(ctx context.Context) {
	// TODO: review this, if we need ctx with differnt span than incoming request , because tracer can work a litter longer then actual request, without affecting the latency of call
	// or have a timeout for this function ...
	for {
		select {
		case <-t.CloseUpdater:
			fmt.Println("Calling Save")
			ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)
			// get timestamp of resp, and populate duration here
			t.db.SaveTrace(ctx,
				t.traceRecord)
			return
		case apply := <-t.Updater:
			apply(t.traceRecord)
		}
	}
}

func (t *TraceRecorder) Record() {
	go t.recorder(context.Background())
}

func NewTraceRecorder(traceID string, db db.TraceStorage) *TraceRecorder {
	out := &TraceRecorder{
		StartTime: time.Now(),
		traceRecord: &models.TraceRecord{
			TraceID:   uuid.NewString(),
			Timestamp: time.Now(),
		},
		Updater:      make(chan func(r *models.TraceRecord)),
		CloseUpdater: make(chan struct{}),
		db:           db,
	}
	inMemoryWriter := &utility.
		InMemoryLogWriter{Writer: out.Updater}
	out.Lg = zerolog.New(inMemoryWriter).With().Timestamp().Logger()
	out.Record()
	return out
}

func (r *TraceRecorder) GetTraceID() string {
	return r.TraceID
}

func FromContext(ctx context.Context) *TraceRecorder {
	traceRecorder, _ := ctx.Value(traceKey).(*TraceRecorder)
	return traceRecorder
}

func ToContext(ctx context.Context,
	traceRecorder *TraceRecorder) context.Context {
	return context.WithValue(ctx, traceKey, traceRecorder)
}
