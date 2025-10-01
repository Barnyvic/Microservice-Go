#!/bin/bash


set -e

echo "=========================================="
echo "Product Microservice Setup"
echo "=========================================="
echo ""

echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.21 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "Go is installed: $GO_VERSION"
echo ""

echo "Checking Protocol Buffers compiler..."
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed."
    echo "   macOS: brew install protobuf"
    echo "   Linux: sudo apt install -y protobuf-compiler"
    echo "   Windows: Download from https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

PROTOC_VERSION=$(protoc --version)
echo " protoc is installed: $PROTOC_VERSION"
echo ""

echo "Installing Go plugins for Protocol Buffers..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo "Go plugins installed"
echo ""

echo "Checking if Go bin is in PATH..."
GOPATH=$(go env GOPATH)
if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
    echo " Warning: $GOPATH/bin is not in your PATH"
    echo "   Add this to your shell profile:"
    echo "   export PATH=\$PATH:\$(go env GOPATH)/bin"
    echo ""
fi

echo "Downloading Go dependencies..."
go mod download
echo "Dependencies downloaded"
echo ""

echo "Generating Protocol Buffer files..."
if command -v make &> /dev/null; then
    make proto
else
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/product.proto proto/subscription.proto
fi
echo "Proto files generated"
echo ""

echo "Running tests..."
go test -v ./... || {
    echo " Some tests failed. This might be expected if the database is not set up."
}
echo ""

echo "Checking for grpcurl (optional but recommended)..."
if command -v grpcurl &> /dev/null; then
    echo "grpcurl is installed"
else
    echo "grpcurl is not installed (optional)"
    echo "   macOS: brew install grpcurl"
    echo "   Linux/Windows: https://github.com/fullstorydev/grpcurl/releases"
fi
echo ""

if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo ".env file created"
else
    echo ".env file already exists"
fi
echo ""

echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Review the .env file and adjust settings if needed"
echo "2. Start the server: make run"
echo "3. In another terminal, test the API:"
echo "   grpcurl -plaintext localhost:50051 list"
echo "4. Or run the example client:"
echo "   go run examples/client/main.go"
echo ""
echo "For more information, see:"
echo "- README.md for detailed documentation"
echo ""

