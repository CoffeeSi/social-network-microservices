package mongo

import (
	"context"

	"github.com/CoffeeSi/social-network-microservices/chat-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	database *mongo.Database
	client   *mongo.Client
}

func NewMongoDatabase(ctx context.Context, cfg *config.Config) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.MongoName)
	return &MongoDB{database: db, client: client}, nil
}

func (m *MongoDB) Database() *mongo.Database {
	return m.database
}
