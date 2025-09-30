@echo off
REM Script to test the gRPC API using grpcurl on Windows
REM Make sure the server is running before executing this script

echo Testing Product Microservice API
echo =================================
echo.

REM Check if grpcurl is installed
where grpcurl >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo grpcurl is not installed. Please install it first:
    echo   Download from: https://github.com/fullstorydev/grpcurl/releases
    exit /b 1
)

set SERVER=localhost:50051

echo 1. Listing available services...
grpcurl -plaintext %SERVER% list
echo.

echo 2. Creating a product...
grpcurl -plaintext -d "{\"name\": \"Premium Software License\", \"description\": \"Enterprise software solution\", \"price\": 299.99, \"product_type\": \"digital\"}" %SERVER% product.ProductService/CreateProduct
echo.

echo 3. Listing all products...
grpcurl -plaintext -d "{\"page\": 1, \"page_size\": 10}" %SERVER% product.ProductService/ListProducts
echo.

echo =================================
echo API Testing Complete!
echo =================================
echo.
echo Note: For full testing with variable capture, use PowerShell or the bash script.

