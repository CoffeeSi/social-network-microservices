package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB(uri string) (*mongo.Client, error) {

	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpt)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}
