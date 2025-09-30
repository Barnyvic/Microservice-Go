#!/bin/bash

# Script to test the gRPC API using grpcurl
# Make sure the server is running before executing this script

set -e

echo "Testing Product Microservice API"
echo "================================="
echo ""

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo "grpcurl is not installed. Please install it first:"
    echo "  macOS: brew install grpcurl"
    echo "  Linux: https://github.com/fullstorydev/grpcurl/releases"
    echo "  Windows: https://github.com/fullstorydev/grpcurl/releases"
    exit 1
fi

SERVER="localhost:50051"

echo "1. Listing available services..."
grpcurl -plaintext $SERVER list
echo ""

echo "2. Creating a product..."
PRODUCT_RESPONSE=$(grpcurl -plaintext -d '{
  "name": "Premium Software License",
  "description": "Enterprise software solution with full support",
  "price": 299.99,
  "product_type": "digital"
}' $SERVER product.ProductService/CreateProduct)

echo "$PRODUCT_RESPONSE"
PRODUCT_ID=$(echo "$PRODUCT_RESPONSE" | grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4)
echo "Created Product ID: $PRODUCT_ID"
echo ""

echo "3. Getting the product..."
grpcurl -plaintext -d "{\"id\": \"$PRODUCT_ID\"}" $SERVER product.ProductService/GetProduct
echo ""

echo "4. Creating a subscription plan..."
PLAN_RESPONSE=$(grpcurl -plaintext -d "{
  \"product_id\": \"$PRODUCT_ID\",
  \"plan_name\": \"Monthly Plan\",
  \"duration\": 30,
  \"price\": 29.99
}" $SERVER subscription.SubscriptionService/CreateSubscriptionPlan)

echo "$PLAN_RESPONSE"
PLAN_ID=$(echo "$PLAN_RESPONSE" | grep -o '"id": "[^"]*"' | head -1 | cut -d'"' -f4)
echo "Created Plan ID: $PLAN_ID"
echo ""

echo "5. Creating another subscription plan..."
grpcurl -plaintext -d "{
  \"product_id\": \"$PRODUCT_ID\",
  \"plan_name\": \"Annual Plan\",
  \"duration\": 365,
  \"price\": 299.99
}" $SERVER subscription.SubscriptionService/CreateSubscriptionPlan
echo ""

echo "6. Listing subscription plans for the product..."
grpcurl -plaintext -d "{\"product_id\": \"$PRODUCT_ID\"}" $SERVER subscription.SubscriptionService/ListSubscriptionPlans
echo ""

echo "7. Updating the product..."
grpcurl -plaintext -d "{
  \"id\": \"$PRODUCT_ID\",
  \"name\": \"Premium Software License - Updated\",
  \"description\": \"Enterprise software solution with full support and updates\",
  \"price\": 349.99,
  \"product_type\": \"digital\"
}" $SERVER product.ProductService/UpdateProduct
echo ""

echo "8. Listing all products..."
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' $SERVER product.ProductService/ListProducts
echo ""

echo "9. Creating a physical product..."
grpcurl -plaintext -d '{
  "name": "Hardware Device",
  "description": "Physical hardware product",
  "price": 499.99,
  "product_type": "physical"
}' $SERVER product.ProductService/CreateProduct
echo ""

echo "10. Listing products filtered by type (digital)..."
grpcurl -plaintext -d '{
  "product_type": "digital",
  "page": 1,
  "page_size": 10
}' $SERVER product.ProductService/ListProducts
echo ""

echo "11. Deleting the subscription plan..."
grpcurl -plaintext -d "{\"id\": \"$PLAN_ID\"}" $SERVER subscription.SubscriptionService/DeleteSubscriptionPlan
echo ""

echo "12. Deleting the product..."
grpcurl -plaintext -d "{\"id\": \"$PRODUCT_ID\"}" $SERVER product.ProductService/DeleteProduct
echo ""

echo "================================="
echo "API Testing Complete!"
echo "================================="

