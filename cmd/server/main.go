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

	dbConfig := database.Config{
		Driver:   getEnv("DB_DRIVER", constants.DefaultDBDriver),
		Host:     getEnv("DB_HOST", constants.DefaultDBHost),
		Port:     getEnv("DB_PORT", constants.DefaultDBPort),
		User:     getEnv("DB_USER", constants.DefaultDBUser),
		Password: getEnv("DB_PASSWORD", constants.DefaultDBPassword),
		DBName:   getEnv("DB_NAME", constants.DefaultDBName),
		SSLMode:  getEnv("DB_SSLMODE", constants.DefaultDBSSLMode),
	}

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("✗ Failed to connect to database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("✗ Failed to run migrations: %v", err)
	}
	
	productRepo := repository.NewProductRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	productService := service.NewProductService(productRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, productRepo)

	productHandler := handler.NewProductHandler(productService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	grpcServer := grpc.NewServer()
		
	productpb.RegisterProductServiceServer(grpcServer, productHandler)
	subscriptionpb.RegisterSubscriptionServiceServer(grpcServer, subscriptionHandler)

	reflection.Register(grpcServer)
	port := getEnv("PORT", constants.DefaultGRPCPort)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("✗ Failed to listen on port %s: %v", port, err)
	}

	log.Println("========================================")
	log.Printf("✓ gRPC server listening on port %s", port)
	log.Println("✓ Server ready to accept connections")
	log.Println("========================================")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("✗ Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\n========================================")
	log.Println("  Shutting down gRPC server...")
	log.Println("========================================")

	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Server gracefully stopped")
	case <-ctx.Done():
		log.Println("Shutdown timeout exceeded, forcing stop")
		grpcServer.Stop()
	}

	log.Println("========================================")
	log.Println("  Server shutdown complete")
	log.Println("========================================")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
