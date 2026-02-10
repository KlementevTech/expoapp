# Makefile

lint:
	golangci-lint run ./...

fix:
	golangci-lint run --fix ./...

GOLANGCI_LINT_VERSION := v2.8.0

install_linter:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

.PHONY: \
	lint \
	fix \
	install_linter