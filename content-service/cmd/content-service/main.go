package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/cache"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/client"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/config"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/events"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository/db"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/transport"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase"
	pb "github.com/IsFariza/maxat-protobuf/content-service-pb"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoDB, err := db.InitMongoDB(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	log.Printf("Connected successfully to MongoDB Database: %s", mongoDB.Name())
	if err := db.RunMigrations(ctx, mongoDB); err != nil {
		log.Fatalf("Failed to run content-service migrations: %v", err)
	}
	log.Println("Content-service MongoDB migrations completed")

	redisAddr := os.Getenv("REDIS_ADDR")

	rClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis instance unreachable. Proceeding with cache bypass: %v", err)
	}

	userSvcURL := os.Getenv("USER_SERVICE_URL")

	userConn, err := grpc.NewClient(userSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userServiceClient := client.NewGrpcUserClient(userConn)

	natsURL := os.Getenv("NATS_URL")

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("unable to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Printf("Connected successfully to NATS at %s", natsURL)
	postPublisher := events.NewPostPublisher(nc)

	postsCollection := mongoDB.Collection("posts")
	likesCollection := mongoDB.Collection("likes")
	commentsCollection := mongoDB.Collection("comments")

	rawPostRepo := repository.NewPostRepo(postsCollection)
	likeRepo := repository.NewLikeRepo(likesCollection)
	commentRepo := repository.NewCommentRepo(commentsCollection)

	redisPostCache := cache.NewRedisPostCache(rClient)
	cachedPostRepo := repository.NewCachedPostRepository(rawPostRepo, redisPostCache)

	postUC := usecase.NewPostUseCase(cachedPostRepo, userServiceClient, postPublisher)
	likeUC := usecase.NewLikeUseCase(likeRepo, cachedPostRepo, userServiceClient)
	commentUC := usecase.NewCommentUseCase(commentRepo, cachedPostRepo, userServiceClient)
	likeQueue := events.NewLikeQueue(nc, likeUC)
	if err := likeQueue.Start(); err != nil {
		log.Fatalf("failed to start NATS: %v", err)
	}
	defer likeQueue.Stop()
	likeUC.SetCommandBus(likeQueue)

	grpcHandler := transport.NewContentHandler(postUC, commentUC, likeUC)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to bind network socket on port %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContentServiceServer(grpcServer, grpcHandler)

	go func() {
		log.Printf("Content gRPC active on port :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("fatal gRPC network execution breakdown: %v", err)
		}
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	log.Println("Shutting down Content microservice gracefully...")
	grpcServer.GracefulStop()
	_ = rClient.Close()
	_ = mongoDB.Client().Disconnect(context.Background())
	log.Println("Content microservice successfully stopped.")
}
