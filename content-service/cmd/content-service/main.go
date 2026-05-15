package main

import (
	"context"
	"log"
	"os"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/cache"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/events"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// 1. Initialize MongoDB (The "Raw" Repository)
	mongoClient := initMongoClient() // Your helper function
	db := mongoClient.Database("content_db")
	rawRepo := mongodb.NewPostRepository(db)

	// 2. Initialize Redis (The Cache Layer)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	// Check connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	postCache := cache.NewRedisPostCache(rdb)

	// 3. Initialize the Proxy (The "Magic" Layer)
	// This wraps the rawRepo with caching logic
	proxyRepo := repository.NewCachedPostRepository(rawRepo, postCache)

	// 4. Initialize NATS (The Event Layer)
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	postPublisher := event.NewPostPublisher(nc)

	// 5. Initialize gRPC Clients (e.g., User Service)
	userClient := initUserServiceClient() // Your helper function

	// 6. Initialize UseCase with the PROXY and PUBLISHER
	postUC := usecase.NewPostUseCase(proxyRepo, userClient, postPublisher)

	// 7. Start the Server (gRPC or HTTP)
	// startServer(postUC)
}
