package mongo

import (
	"context"
	"log"

	"time"

	"github.com/p4-pankaj/trace-replay/config"
	"github.com/p4-pankaj/trace-replay/constants"
	"github.com/p4-pankaj/trace-replay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const traceDbName = "myAppDatabase"
const traceCollectionName = "traceCol"

func NewMongoDb(conf *config.MongoConfig) (db *MongoDatabase, err error) {
	clientOptions := options.Client().ApplyURI(conf.URI)
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(conf.ConnectTimeout)*time.Second)
	defer cancel()

	var client *mongo.Client
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("%s: failed to connect to MongoDB: %v",
			constants.PkgTitle, err)
		return
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("%s: failed to ping MongoDB: %v",
			constants.PkgTitle, err)
		return
	}
	db = &MongoDatabase{client: client}
	log.Printf("%s: connected to MongoDB!", constants.PkgTitle)
	return
}

type MongoDatabase struct {
	client *mongo.Client
}

func (m *MongoDatabase) SaveTrace(ctx context.Context,
	trace *models.TraceRecord) error {
	collection := m.client.Database(traceDbName).
		Collection(traceCollectionName)

	// Insert trace data into MongoDB collection
	_, err := collection.InsertOne(ctx, trace)
	if err != nil {
		log.Printf("%s: failed to insert trace data: %v",
			constants.PkgTitle, err)
		return err
	}

	log.Printf("%s: successfully inserted trace data",
		constants.PkgTitle)
	return nil
}

func (m *MongoDatabase) GetTraceByID(ctx context.Context,
	traceID string) (*models.TraceRecord, error) {
	collection := m.client.Database(traceDbName).Collection(traceCollectionName)

	filter := bson.D{{Key: "_id", Value: traceID}}

	var trace models.TraceRecord
	err := collection.FindOne(ctx, filter).Decode(&trace)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("%s: trace not found for ID: %s",
				constants.PkgTitle, traceID)
			return nil, nil // Return nil if not found
		}
		log.Printf("%s: failed to retrieve trace for ID %s: %v",
			constants.PkgTitle, traceID, err)
		return nil, err
	}

	log.Printf("%s: successfully retrieved trace for ID %s",
		constants.PkgTitle, traceID)
	return &trace, nil
}
