# Makefile

lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...

PROTO_DIR := api/expo
OUT_DIR := pkg/api/expo

gen_protobuf:
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

GOLANGCI_LINT_VERSION := v2.8.0

install_linter:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

install_protoc:
	@echo "Installing protoc..."
	sudo apt install -y protobuf-compiler

install_tools: install_linter
	@echo "Installing proto plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: \
	lint \
	fix \
	gen_protobuf \
	install_linter \
	install_protoc \
	install_tools