package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(uri string) (*mongo.Client, error) {

	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOpt)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}
