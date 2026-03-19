MODULE := $(shell go list -m)

PROTO_DIR := api
OUT_DIR := .

PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

.PHONY: all generate clean help install-deps

all: generate

generate:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)
	@echo "✅ Kod wygenerowany pomyślnie!"

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Wtyczki zainstalowane. Upewnij się, że $(shell go env GOPATH)/bin jest w Twoim PATH."

clean:
	find . -name "*.pb.go" -type f -delete
	@echo "🗑️ Wygenerowane pliki usunięte."

help:
	@echo "Dostępne komendy:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'