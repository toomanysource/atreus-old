##@ Proto Development

API_PROTO_FILES := $(shell find api -name *.proto)
APP_PROTO_FILES := $(shell find app -name *.proto)

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: proto
proto: ## Generate code based on the api and app proto files.
proto: proto.gen.all

.PHONY: proto-api
proto-api: ## Generate code based on the api proto files.
proto-api: proto.gen.api

.PHONY: proto-app
proto-app: ## Generate code based on the app proto files.
proto-app: proto.gen.app

# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: proto.gen.api
proto.gen.api:
	protoc --proto_path=./api \
    	   --proto_path=./third_party \
     	   --go_out=paths=source_relative:./api \
     	   --go-http_out=paths=source_relative:./api \
     	   --go-grpc_out=paths=source_relative:./api \
     	   --validate_out=paths=source_relative,lang=go:./api \
  		   $(API_PROTO_FILES)

.PHONY: proto.gen.app
proto.gen.app:
	protoc --proto_path=./app \
    	   --proto_path=./third_party \
     	   --go_out=paths=source_relative:./app \
    	   $(APP_PROTO_FILES)

.PHONY: proto.gen.all
proto.gen.all: proto.gen.api proto.gen.app
