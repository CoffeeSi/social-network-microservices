package app

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/CoffeeSi/social-network-microservices/user-service/internal/repository"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/repository/db"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/transport"
	"github.com/CoffeeSi/social-network-microservices/user-service/internal/usecase"
	pb "github.com/CoffeeSi/social-network-microservices/user-service/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *Config) {
	client, err := db.InitMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("fail to connect to mongo %v", err)
	}
	db := client.Database("user_db")
	//mq, err := event.NewNatsPublisher(cfg.NatsURI)
	// if err != nil {
	// 	log.Printf("Failed to connect to NATS: $v", err)
	// }
	opt, err := redis.ParseURL(cfg.REDIS)
	if err != nil {
		log.Printf("failed to parse redis url, using default: %v", err)
		opt = &redis.Options{Addr: "localhost:6379"}
	}
	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: redis unavailable, caching disabled: %v", err)
	}

	userRepo := repository.NewUserRepo(db.Collection("users"))
	ttl, _ := strconv.Atoi(cfg.TTL)
	//rpm, _ := strconv.Atoi(cfg.RPM)
	cachedRepo := repository.NewCachedUserRepository(userRepo, rdb, time.Duration(ttl)*time.Second)
	// rateLimiter := middleware.NewRateLimiter(rdb, rpm)
	userService := usecase.NewUserUseCase(cachedRepo) //mq later
	userHandler := transport.NewUserHandler(userService)

	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
	//grpc.UnaryInterceptor(rateLimiter.UnaryInterceptor()),
	)
	pb.RegisterUserServiceServer(s, userHandler)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
