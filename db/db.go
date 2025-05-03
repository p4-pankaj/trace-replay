package db

import (
	"context"

	"github.com/p4-pankaj/trace-replay/db/mongo.go"
	"github.com/p4-pankaj/trace-replay/models"
	"github.com/p4-pankaj/trace-replay/traceConfig"
)

// TraceStorage defines the interface for storing and retrieving trace data.
// Users can pass either MySQL or MongoDB config for trace storage.
// If they wish to use another database, they can provide an interface
// implementing these methods.
type TraceStorage interface {
	SaveTrace(ctx context.Context,
		trace *models.TraceRecord) error

	GetTraceByID(ctx context.Context,
		traceID string) (*models.TraceRecord, error)
}

func InitTraceStorage(c *traceConfig.DbConfig) (store TraceStorage,
	err error) {
	if c.DbKind == traceConfig.MongoDbType {
		return mongo.NewMongoDb(c.MongoConfig)
	}
	return
}
