package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CoffeeSi/social-network-microservices/content-service/internal/repository"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/transport"
	"github.com/CoffeeSi/social-network-microservices/content-service/internal/usecase"
	pb "github.com/CoffeeSi/social-network-microservices/content-service/proto"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockUserServiceClient struct{}

func (m *MockUserServiceClient) UserExists(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

func main() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := "mongodb://localhost:27017"
	log.Printf("Connecting to MongoDB at %s...", mongoURI)

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	db := client.Database("social_network_content_test")

	log.Println("Initializing repositories and usecases...")
	postRepo := repository.NewPostRepo(db.Collection("posts"))
	commentRepo := repository.NewCommentRepo(db.Collection("comments"))

	mockUserClient := &MockUserServiceClient{}

	postUC := usecase.NewPostUseCase(postRepo, mockUserClient)
	commentUC := usecase.NewCommentUseCase(commentRepo, postRepo, mockUserClient)

	handler := transport.NewContentHandler(postUC, commentUC, nil)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContentServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("Content Service is running on gRPC port :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gRPC server gracefully...")
	grpcServer.GracefulStop()
	log.Println("Server stopped.")
}
