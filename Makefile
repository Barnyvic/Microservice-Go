.PHONY: proto build run test clean


proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product.proto proto/subscription.proto

build: proto
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

test:
	go test -v ./...

test-with-cgo:
	CGO_ENABLED=1 go test -v ./...

test-coverage:
	CGO_ENABLED=1 go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f proto/*.pb.go
	rm -f coverage.out coverage.html

deps:
	go mod download
	go mod tidy

install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

