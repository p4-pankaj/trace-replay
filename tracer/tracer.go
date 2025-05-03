package tracer

import (
	"context"
	"fmt"

	"github.com/p4-pankaj/trace-replay/db"
	"github.com/p4-pankaj/trace-replay/models"
	"github.com/p4-pankaj/trace-replay/utility"

	"time"

	"github.com/rs/zerolog"
)

type contextKey string

const traceKey contextKey = "traceRecorder"

type TraceRecorder struct {
	Env          string
	TraceID      string
	StartTime    time.Time
	TraceRecord  *models.TraceRecord
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
			ct, cancel := context.WithTimeout(context.Background(),
				15*time.Second)
			defer cancel()
			// get timestamp of resp, and populate duration here

			err := t.db.SaveTrace(ct,
				t.TraceRecord)
			// Gracefully bypassing trace-level error for now.
			// In future versions, users may configure an environment variable to treat this as a critical error.
			if err != nil {
				fmt.Println("error saving trace", err)
				return
			}
		case apply := <-t.Updater:
			apply(t.TraceRecord)
		}
	}
}

func (t *TraceRecorder) Record() {
	go t.recorder(context.Background())
}

func NewTraceRecorder(ctx context.Context,
	traceID string,
	db db.TraceStorage,
	env string,
) (out *TraceRecorder) {

	if env == "DEBUG" {
		traceRecord, err := db.GetTraceByID(ctx, traceID)
		if err != nil {
			panic("failed to fetch trace for this traceId")
		}
		out = &TraceRecorder{
			Env:         env,
			StartTime:   time.Now(),
			TraceRecord: traceRecord,
		}
	} else {
		out = &TraceRecorder{
			Env:       env,
			StartTime: time.Now(),
			TraceRecord: &models.TraceRecord{
				TraceID:   traceID,
				Timestamp: time.Now(),
			},
			Updater:      make(chan func(r *models.TraceRecord)),
			CloseUpdater: make(chan struct{}),
			db:           db,
		}
		out.Record()
		inMemoryWriter := &utility.
			InMemoryLogWriter{Writer: out.Updater}
		out.Lg = zerolog.New(inMemoryWriter).With().Timestamp().Logger()
	}
	return out
}

func (r *TraceRecorder) GetTraceID() string {
	return r.TraceID
}

// Wrapper function that wraps a caller function, taking an argument of type T and returning a value of type Q
func Wrapper[T any, Q any](r *TraceRecorder,
	caller func(field T) Q,
	field2 T) (out Q) {
	// Call the function with the provided field and return the result
	if r.Env == "DEBUG" {

	} else {
		if r.TraceRecord.FunctionTrace == nil {
			r.TraceRecord.FunctionTrace = []models.FunctionTrace{}
		}
		funcTrace := models.FunctionTrace{}
		funcTrace.Input = field2
		funcTrace.ReqTimestamp = time.Now()
		funcTrace.CallID = utility.GetHashForObject(field2)

		out = caller(field2)

		funcTrace.Output = out
		funcTrace.RespTimestamp = time.Now()
		funcTrace.Duration = funcTrace.RespTimestamp.
			Sub(funcTrace.ReqTimestamp)

		r.TraceRecord.FunctionTrace = append(r.TraceRecord.FunctionTrace, funcTrace)
		r.Updater <- func(t *models.TraceRecord) {
			t.FunctionTrace = r.TraceRecord.FunctionTrace
		}

	}
	return out
}

func FromContext(ctx context.Context) *TraceRecorder {
	traceRecorder, _ := ctx.Value(traceKey).(*TraceRecorder)
	return traceRecorder
}

func ToContext(ctx context.Context,
	traceRecorder *TraceRecorder) context.Context {
	return context.WithValue(ctx, traceKey, traceRecorder)
}
