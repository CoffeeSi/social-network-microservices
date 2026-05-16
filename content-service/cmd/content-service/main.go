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
	pb "github.com/CoffeeSi/social-network-microservices/content-service/proto"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Initializing Maxat Content Microservice...")

	// 1. Load system configuration environment variables
	cfg := config.LoadConfig()

	// 2. Establish connection to MongoDB cluster passing context parameter
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoDB, err := db.InitMongoDB(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Critical failure connecting to MongoDB: %v", err)
	}
	log.Printf("Connected successfully to MongoDB Database: %s", mongoDB.Name())

	// 3. Connect to Redis memory cache cluster
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if err := rClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("Warning: Redis instance unreachable. Proceeding with cache bypass: %v", err)
	}

	// 4. Initialize User Service downstream gRPC Client Connection
	userSvcURL := os.Getenv("USER_SERVICE_URL")
	if userSvcURL == "" {
		userSvcURL = "localhost:50052"
	}
	userConn, err := grpc.Dial(userSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to establish gRPC dial link to User Service: %v", err)
	}
	defer userConn.Close()
	userServiceClient := client.NewGrpcUserClient(userConn)

	// 5. Initialize NATS Post Event Publisher
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Critical error: unable to establish connection to NATS message broker: %v", err)
	}
	defer nc.Close()
	log.Printf("Connected successfully to NATS Message Broker at %s", natsURL)
	postPublisher := events.NewPostPublisher(nc)

	// 6. Initialize MongoDB collections
	postsCollection := mongoDB.Collection("posts")
	likesCollection := mongoDB.Collection("likes")
	commentsCollection := mongoDB.Collection("comments")

	// 7. Wire up concrete infrastructure repositories
	rawPostRepo := repository.NewPostRepo(postsCollection)
	likeRepo := repository.NewLikeRepo(likesCollection)
	commentRepo := repository.NewCommentRepo(commentsCollection)

	// 8. Wrap the post repository into your optimized Global Cache Proxy
	redisPostCache := cache.NewRedisPostCache(rClient)
	cachedPostRepo := repository.NewCachedPostRepository(rawPostRepo, redisPostCache)

	// 9. Instantiate business orchestration usecases with all required dependencies
	postUC := usecase.NewPostUseCase(cachedPostRepo, userServiceClient, postPublisher)
	likeUC := usecase.NewLikeUseCase(likeRepo, cachedPostRepo, userServiceClient)
	commentUC := usecase.NewCommentUseCase(commentRepo, cachedPostRepo, userServiceClient)

	// 10. Wire presentation layer dependencies into unified gRPC handler
	grpcHandler := transport.NewContentHandler(postUC, commentUC, likeUC)

	// 11. Bind microservice onto the target network TCP port socket
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to bind network socket on port %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContentServiceServer(grpcServer, grpcHandler)

	// 12. Start listening gRPC engine inside a non-blocking background thread
	go func() {
		log.Printf("Content microservice gRPC engine active on port :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Fatal gRPC network execution breakdown: %v", err)
		}
	}()

	// 13. Graceful shutdown block tracking OS container stop triggers
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	log.Println("Shutting down Content microservice containers gracefully...")
	grpcServer.GracefulStop()
	_ = rClient.Close()
	_ = mongoDB.Client().Disconnect(context.Background())
	log.Println("Content Microservice successfully offline.")
}
