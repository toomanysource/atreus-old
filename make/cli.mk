##@ Cli Development

WIRE := wire
WIRE_DIRS := $(abspath $(dir $(shell find -name wire.go)))

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: wire
wire: ## Generate wire_gen code based on every wire.go
wire: wire.gen

.PHONY: cli-install
cli-install: ## Install the project basic CLI tools
cli-install: cli.install


# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: cli.install
cli.install:
	@echo "===========> Installing kratos cli"
	@go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	@echo "===========> Installing wire"
	@go install github.com/google/wire/cmd/wire@latest
	@echo "===========> Installing protoc-gen-go"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo "===========> Installing protoc-gen-go-grpc"
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "===========> Installing protoc-gen-go-http"
	@go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest

.PHONY: wire.gen
wire.gen:
	@$(WIRE) gen $(WIRE_DIRS)