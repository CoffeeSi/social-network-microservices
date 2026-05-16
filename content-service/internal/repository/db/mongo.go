package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB(ctx context.Context, uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb cluster: %w", err)

	}

	return client.Database(dbName), nil
}
