GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)