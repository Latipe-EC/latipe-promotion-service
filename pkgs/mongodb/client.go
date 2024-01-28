package mongodb

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"latipe-promotion-services/config"
)

type MongoClient struct {
	client  *mongo.Client
	rootCtx context.Context
	cfg     *config.Config
}

// Open - creates a new Mongo
func OpenMongoDBConnection(cfg *config.Config) (*MongoClient, error) {
	ctx := context.Background()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			fmt.Println(e.Command)
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			fmt.Println(failedEvent.Failure)
		},
	}

	opts := options.Client().SetMonitor(monitor)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongodb.ConnectionString), opts)
	if err != nil {
		return nil, err
	}

	return &MongoClient{client: client, rootCtx: ctx, cfg: cfg}, nil

}

func (m *MongoClient) GetDB() *mongo.Database {
	db := m.client.Database(m.cfg.Mongodb.DbName)
	return db
}

// Disconnect - used mainly in testing to avoid capping out the concurrent connections on MongoDB
func (m *MongoClient) Disconnect() {
	err := m.client.Disconnect(m.rootCtx)
	if err != nil {
		log.Fatalf("disconnecting from mongodb: %v", err)
	}
}

// Ping sends a ping command to verify that the client can connect to the deployment.
func (m *MongoClient) Ping() error {
	return m.client.Ping(m.rootCtx, readpref.Primary())
}
