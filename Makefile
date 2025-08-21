GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: generate-openapi
# generate OpenAPI documentation
generate-openapi:
	@echo "Generating OpenAPI documentation..."
	mkdir -p docs/api/openapi
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --openapi_out=fq_schema_naming=true,default_response=false,allow_merge=true:./docs/api/openapi \
	       $(API_PROTO_FILES)
	@echo "OpenAPI documentation generated in docs/api/openapi/"

.PHONY: generate-swagger-ui
# generate Swagger UI
generate-swagger-ui:
	@echo "Setting up Swagger UI..."
	mkdir -p docs/api/swagger
	@if [ ! -d "docs/api/swagger/swagger-ui" ]; then \
		curl -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v5.10.3.tar.gz | tar -xz; \
		mv swagger-ui-5.10.3/dist docs/api/swagger/swagger-ui; \
		rm -rf swagger-ui-5.10.3; \
	fi
	@echo "Swagger UI setup completed"

.PHONY: generate-docs
# generate all documentation
generate-docs: generate-openapi generate-swagger-ui
	@echo "All documentation generated successfully"

.PHONY: serve-docs
# serve documentation locally
serve-docs:
	@echo "Starting documentation server at http://localhost:8080/docs"
	@cd docs/api && python3 -m http.server 8080

.PHONY: watch-docs
# watch for changes and regenerate docs
watch-docs:
	@echo "Watching for .proto file changes..."
	@while true; do \
		inotifywait -r -e modify,create,delete api/; \
		make generate-openapi; \
		echo "Documentation updated at $$(date)"; \
	done

.PHONY: build
# build
build: generate-docs
	@echo "Building with version: $(VERSION)"
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;
	make generate-docs;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
