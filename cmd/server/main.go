package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/microservice-go/product-service/internal/constants"
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
	log.Println("========================================")
	log.Println("  Product Microservice Starting...")
	log.Println("========================================")

	// Database configuration
	dbConfig := database.Config{
		Driver:   getEnv("DB_DRIVER", constants.DefaultDBDriver),
		Host:     getEnv("DB_HOST", constants.DefaultDBHost),
		Port:     getEnv("DB_PORT", constants.DefaultDBPort),
		User:     getEnv("DB_USER", constants.DefaultDBUser),
		Password: getEnv("DB_PASSWORD", constants.DefaultDBPassword),
		DBName:   getEnv("DB_NAME", constants.DefaultDBName),
		SSLMode:  getEnv("DB_SSLMODE", constants.DefaultDBSSLMode),
	}

	// Initialize database
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("✗ Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("✗ Failed to run migrations: %v", err)
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
	port := getEnv("PORT", constants.DefaultGRPCPort)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("✗ Failed to listen on port %s: %v", port, err)
	}

	log.Println("========================================")
	log.Printf("✓ gRPC server listening on port %s", port)
	log.Println("✓ Server ready to accept connections")
	log.Println("========================================")

	// Setup graceful shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("✗ Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n========================================")
	log.Println("  Shutting down gRPC server...")
	log.Println("========================================")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("✓ Server gracefully stopped")
	case <-ctx.Done():
		log.Println("⚠ Shutdown timeout exceeded, forcing stop")
		grpcServer.Stop()
	}

	log.Println("========================================")
	log.Println("  Server shutdown complete")
	log.Println("========================================")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
