include .env
export

NAMESPACE := klementevtech
APP_NAME := expoapp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")
DOCKER_TAG := $(NAMESPACE)/$(APP_NAME):$(VERSION)

lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...

test:
	@echo "Running unit tests..."
	go test ./...

build_bin:
	@echo "Building bin file, version [$(VERSION)]..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=$(VERSION) -s -w" -o bin/$(APP_NAME) ./cmd/server

docker_build: build_bin
	@echo "Building docker image, tag [$(DOCKER_TAG)]..."
	docker build -t $(DOCKER_TAG) -f ./ci/docker/Dockerfile .

up:
	VERSION=$(VERSION) docker compose -p expoapp -f ./ci/dev/docker-compose.yml --env-file .env up -d

PROTO_DIR := api/proto
OUT_DIR := pkg/api/pb

gen_pb:
	@echo "Creating dir for protobuf files..."
	mkdir -p $(OUT_DIR)
	@echo "Generating protobuf files..."
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

gen_mocks:
	go generate generate.go

gen_all: gen_pb gen_mocks

GOLANGCI_LINT_VERSION := v2.9.0

install_linter:
	@echo "Installing golangci-lint, version [$(GOLANGCI_LINT_VERSION)]..."
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

install_protoc:
	@echo "Installing protoc..."
	sudo apt install -y protobuf-compiler

install_tools:
	@echo "Installing proto plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Installing mockgen plugin..."
	go install go.uber.org/mock/mockgen@latest
	@echo "Installing grpcurl..."
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

help:
	cat Makefile