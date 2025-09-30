package main

import (
	"context"
	"fmt"
	"log"
	"time"

	productpb "github.com/microservice-go/product-service/proto/product"
	subscriptionpb "github.com/microservice-go/product-service/proto/subscription"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	productClient := productpb.NewProductServiceClient(conn)
	subscriptionClient := subscriptionpb.NewSubscriptionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== Product Microservice Client Example ===")

	// 1. Create a product
	fmt.Println("1. Creating a product...")
	createProductResp, err := productClient.CreateProduct(ctx, &productpb.CreateProductRequest{
		Name:        "Premium Software License",
		Description: "Enterprise software solution with full support",
		Price:       299.99,
		ProductType: "digital",
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	productID := createProductResp.Product.Id
	fmt.Printf("✓ Created product: %s (ID: %s)\n\n", createProductResp.Product.Name, productID)

	// 2. Get the product
	fmt.Println("2. Retrieving the product...")
	getProductResp, err := productClient.GetProduct(ctx, &productpb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		log.Fatalf("Failed to get product: %v", err)
	}
	fmt.Printf("✓ Retrieved product: %s - $%.2f\n", getProductResp.Product.Name, getProductResp.Product.Price)
	fmt.Printf("  Description: %s\n", getProductResp.Product.Description)
	fmt.Printf("  Type: %s\n\n", getProductResp.Product.ProductType)

	// 3. Create subscription plans
	fmt.Println("3. Creating subscription plans...")
	
	// Monthly plan
	monthlyPlanResp, err := subscriptionClient.CreateSubscriptionPlan(ctx, &subscriptionpb.CreateSubscriptionPlanRequest{
		ProductId: productID,
		PlanName:  "Monthly Plan",
		Duration:  30,
		Price:     29.99,
	})
	if err != nil {
		log.Fatalf("Failed to create monthly plan: %v", err)
	}
	fmt.Printf("✓ Created plan: %s - $%.2f for %d days\n", 
		monthlyPlanResp.Plan.PlanName, monthlyPlanResp.Plan.Price, monthlyPlanResp.Plan.Duration)

	// Annual plan
	annualPlanResp, err := subscriptionClient.CreateSubscriptionPlan(ctx, &subscriptionpb.CreateSubscriptionPlanRequest{
		ProductId: productID,
		PlanName:  "Annual Plan",
		Duration:  365,
		Price:     299.99,
	})
	if err != nil {
		log.Fatalf("Failed to create annual plan: %v", err)
	}
	fmt.Printf("✓ Created plan: %s - $%.2f for %d days\n\n", 
		annualPlanResp.Plan.PlanName, annualPlanResp.Plan.Price, annualPlanResp.Plan.Duration)

	// 4. List subscription plans for the product
	fmt.Println("4. Listing all subscription plans for the product...")
	listPlansResp, err := subscriptionClient.ListSubscriptionPlans(ctx, &subscriptionpb.ListSubscriptionPlansRequest{
		ProductId: productID,
	})
	if err != nil {
		log.Fatalf("Failed to list plans: %v", err)
	}
	fmt.Printf("✓ Found %d subscription plans:\n", listPlansResp.Total)
	for i, plan := range listPlansResp.Plans {
		fmt.Printf("  %d. %s - $%.2f (%d days)\n", i+1, plan.PlanName, plan.Price, plan.Duration)
	}
	fmt.Println()

	// 5. Update the product
	fmt.Println("5. Updating the product...")
	updateProductResp, err := productClient.UpdateProduct(ctx, &productpb.UpdateProductRequest{
		Id:          productID,
		Name:        "Premium Software License - Enterprise Edition",
		Description: "Enterprise software solution with full support, updates, and priority assistance",
		Price:       349.99,
		ProductType: "digital",
	})
	if err != nil {
		log.Fatalf("Failed to update product: %v", err)
	}
	fmt.Printf("✓ Updated product: %s - $%.2f\n\n", updateProductResp.Product.Name, updateProductResp.Product.Price)

	// 6. Create another product
	fmt.Println("6. Creating another product...")
	_, err = productClient.CreateProduct(ctx, &productpb.CreateProductRequest{
		Name:        "Hardware Device",
		Description: "Physical hardware product with warranty",
		Price:       499.99,
		ProductType: "physical",
	})
	if err != nil {
		log.Fatalf("Failed to create second product: %v", err)
	}
	fmt.Println("✓ Created second product")

	// 7. List all products
	fmt.Println("7. Listing all products...")
	listProductsResp, err := productClient.ListProducts(ctx, &productpb.ListProductsRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}
	fmt.Printf("✓ Found %d total products:\n", listProductsResp.Total)
	for i, product := range listProductsResp.Products {
		fmt.Printf("  %d. %s - $%.2f [%s]\n", i+1, product.Name, product.Price, product.ProductType)
	}
	fmt.Println()

	// 8. List products filtered by type
	fmt.Println("8. Listing digital products only...")
	listDigitalResp, err := productClient.ListProducts(ctx, &productpb.ListProductsRequest{
		ProductType: "digital",
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		log.Fatalf("Failed to list digital products: %v", err)
	}
	fmt.Printf("✓ Found %d digital products:\n", listDigitalResp.Total)
	for i, product := range listDigitalResp.Products {
		fmt.Printf("  %d. %s - $%.2f\n", i+1, product.Name, product.Price)
	}
	fmt.Println()

	// 9. Update a subscription plan
	fmt.Println("9. Updating the monthly plan...")
	updatePlanResp, err := subscriptionClient.UpdateSubscriptionPlan(ctx, &subscriptionpb.UpdateSubscriptionPlanRequest{
		Id:        monthlyPlanResp.Plan.Id,
		ProductId: productID,
		PlanName:  "Monthly Plan - Special Offer",
		Duration:  30,
		Price:     24.99,
	})
	if err != nil {
		log.Fatalf("Failed to update plan: %v", err)
	}
	fmt.Printf("✓ Updated plan: %s - $%.2f\n\n", updatePlanResp.Plan.PlanName, updatePlanResp.Plan.Price)

	// 10. Get a specific subscription plan
	fmt.Println("10. Retrieving the updated plan...")
	getPlanResp, err := subscriptionClient.GetSubscriptionPlan(ctx, &subscriptionpb.GetSubscriptionPlanRequest{
		Id: monthlyPlanResp.Plan.Id,
	})
	if err != nil {
		log.Fatalf("Failed to get plan: %v", err)
	}
	fmt.Printf("✓ Retrieved plan: %s - $%.2f for %d days\n\n", 
		getPlanResp.Plan.PlanName, getPlanResp.Plan.Price, getPlanResp.Plan.Duration)

	// 11. Delete a subscription plan
	fmt.Println("11. Deleting the annual plan...")
	deletePlanResp, err := subscriptionClient.DeleteSubscriptionPlan(ctx, &subscriptionpb.DeleteSubscriptionPlanRequest{
		Id: annualPlanResp.Plan.Id,
	})
	if err != nil {
		log.Fatalf("Failed to delete plan: %v", err)
	}
	fmt.Printf("✓ %s\n\n", deletePlanResp.Message)

	// 12. Verify deletion
	fmt.Println("12. Verifying plan deletion...")
	listPlansAfterDelete, err := subscriptionClient.ListSubscriptionPlans(ctx, &subscriptionpb.ListSubscriptionPlansRequest{
		ProductId: productID,
	})
	if err != nil {
		log.Fatalf("Failed to list plans: %v", err)
	}
	fmt.Printf("✓ Now %d subscription plans remain\n\n", listPlansAfterDelete.Total)

	// 13. Delete a product (cascade delete will remove associated plans)
	fmt.Println("13. Deleting the product (cascade delete)...")
	deleteProductResp, err := productClient.DeleteProduct(ctx, &productpb.DeleteProductRequest{
		Id: productID,
	})
	if err != nil {
		log.Fatalf("Failed to delete product: %v", err)
	}
	fmt.Printf("✓ %s\n\n", deleteProductResp.Message)

	fmt.Println("=== All operations completed successfully! ===")
}

