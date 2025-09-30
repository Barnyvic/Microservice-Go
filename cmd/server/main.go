package main

import (
	"log"
	"net"
	"os"

	"github.com/microservice-go/product-service/internal/database"
	"github.com/microservice-go/product-service/internal/handler"
	"github.com/microservice-go/product-service/internal/repository"
	"github.com/microservice-go/product-service/internal/service"
	productpb "github.com/microservice-go/product-service/proto/product"
	subscriptionpb "github.com/microservice-go/product-service/proto/subscription"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Database configuration
	dbConfig := database.Config{
		Driver:   getEnv("DB_DRIVER", "sqlite"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "products.db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, productRepo)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	productpb.RegisterProductServiceServer(grpcServer, productHandler)
	subscriptionpb.RegisterSubscriptionServiceServer(grpcServer, subscriptionHandler)

	// Register reflection service for grpcurl
	reflection.Register(grpcServer)

	// Start server
	port := getEnv("PORT", "50051")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

