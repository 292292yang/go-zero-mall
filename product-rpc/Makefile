# ==============================
# go-zero rpc Makefile
# ==============================

SERVICE ?= $(notdir $(CURDIR))
PROTO   ?= $(firstword $(wildcard proto/*.proto))
MAIN    ?= $(SERVICE).go
CONF    ?= etc/$(SERVICE).yaml

BIN_DIR ?= bin
BIN     ?= $(BIN_DIR)/$(SERVICE)-rpc

GO      ?= go
GOCTL   ?= goctl

.PHONY: help deps proto tidy fmt vet test build run clean

help:
	@echo "Available targets:"
	@echo "  make deps     Install required tools"
	@echo "  make proto    Generate rpc code from proto"
	@echo "  make tidy     Run go mod tidy"
	@echo "  make fmt      Format go files"
	@echo "  make vet      Run go vet"
	@echo "  make test     Run tests"
	@echo "  make build    Build rpc service"
	@echo "  make run      Run rpc service"
	@echo "  make clean    Clean build files"

deps:
	$(GO) install github.com/zeromicro/go-zero/tools/goctl@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto:
	@if [ -z "$(PROTO)" ]; then \
		echo "No .proto file found"; \
		exit 1; \
	fi
	$(GOCTL) rpc protoc $(PROTO) \
		--go_out=. \
		--go-grpc_out=. \
		--zrpc_out=.

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(MAIN)

run:
	$(GO) run $(MAIN) -f $(CONF)

clean:
	rm -rf $(BIN_DIR)