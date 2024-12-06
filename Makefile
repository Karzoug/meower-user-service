# ==============================================================================
# Define dependencies

SERVICE_VERSION            := 0.1.0
SERVICE_NAME               := user-service
OUTBOX_VERSION             := 0.1.0
BUILD_VERSION              ?= $(shell git symbolic-ref HEAD 2> /dev/null | cut -b 12-)-$(shell git log --pretty=format:\"%h\" -1)
BUILD_DATE                 ?= $(shell date +%FT%T%z)

BASE_IMAGE                 := meower
IMAGE_SERVICE              := $(BASE_IMAGE)/service/user:$(SERVICE_VERSION)
IMAGE_OUTBOX               := $(BASE_IMAGE)/outbox/user:$(OUTBOX_VERSION)

GOLANGCI_LINT_VERSION      := 1.61.0
BUF_VERSION                := 1.46.0
PROTOC_GEN_GO_VERSION 	   := 1.35.1
PROTOC_GEN_GO_GRPC_VERSION := 1.5.1
MIGRATE_VERSION            := 4.18.1

MAIN_PACKAGE_PATH          := ./cmd/
BINARY_NAME                := user_service
TEMP_DIR                   := /var/tmp/meower/user
TEMP_BIN                   := ${TEMP_DIR}/bin
PROJECT_PKG                := github.com/Karzoug/meower-user-service

LDFLAGS += -s -w
LDFLAGS += -X ${PROJECT_PKG}/pkg/buildinfo.buildVersion=${BUILD_VERSION} -X ${PROJECT_PKG}/pkg/buildinfo.buildDate=${BUILD_DATE} -X ${PROJECT_PKG}/pkg/buildinfo.serviceVersion=$(SERVICE_VERSION)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test fmt lint
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## fmt: format .go files
.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -local=${PROJECT_PKG} -l -w  .
	go run mvdan.cc/gofumpt@latest -l -w  .

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## lint: run linters
.PHONY: lint
lint:
	$(TEMP_BIN)/golangci-lint run ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	go build -ldflags "${LDFLAGS}" -o ${TEMP_BIN}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## generate: generate all necessary code
.PHONY: generate
generate:
	$(TEMP_BIN)/buf generate --template buf.gen.grpc.yaml
	$(TEMP_BIN)/buf generate --template buf.gen.kafka.delivery.yaml

## clean: clean all temporary files
.PHONY: clean
clean:
	rm -rf $(TEMP_DIR)

# ==============================================================================
# Database migrations

migrate-create:
	$(TEMP_BIN)/migrate create -ext .sql -dir migrations $(name)

# ==============================================================================
# Kafka tests

## dev-produce-user: send a message to Kafka - user created (ex: dev-produce-user username="bob")
dev-produce-user:
	mkdir ${TEMP_DIR}/proto && cd ${TEMP_DIR}/proto && \
	curl -H 'Accept: application/vnd.github.v3.raw' -O -L https://api.github.com/repos/Karzoug/meower-api/contents/proto/auth/v1/kafka.proto
	$(TEMP_BIN)/protokaf produce ChangedEvent \
		--broker localhost:9094 \
		--proto ${TEMP_DIR}/proto/kafka.proto \
		--topic auth \
		--header "fngpnt=cee25ed946d449a7fb6ded06756619ba" \
		--data '{"username": "$(username)", "change_type": 1}'
	rm -rf ${TEMP_DIR}/proto

# ==============================================================================
# Install dependencies

## dev-install-deps: install dependencies with fixed versions in a temporary directory
dev-install-deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TEMP_BIN) v${GOLANGCI_LINT_VERSION}
	GOBIN=$(TEMP_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}
	GOBIN=$(TEMP_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VERSION}
	GOBIN=$(TEMP_BIN) go install github.com/bufbuild/buf/cmd/buf@v$(BUF_VERSION)
	GOBIN=$(TEMP_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v$(MIGRATE_VERSION)
	GOBIN=$(TEMP_BIN) go install github.com/kuper-tech/protokaf@latest

# ==============================================================================
# Building containers

## service: build the service image
.PHONY: service
service:
	docker build \
		-f build/dockerfile.service \
		-t $(IMAGE_SERVICE) \
		--build-arg BUILD_REF=$(BUILD_VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg VERSION=$(SERVICE_VERSION) \
		--build-arg PROJECT_PKG=$(PROJECT_PKG) \
		--build-arg SERVICE_NAME=$(SERVICE_NAME) \
		.

# ------------------------------------------------------------------------------
# Run docker compose to test service locally

dev-compose-up:
	IMAGE_SERVICE=$(IMAGE_SERVICE) IMAGE_OUTBOX=$(IMAGE_OUTBOX) docker compose -f deploy/docker/dev-compose.yaml  up -d

dev-compose-down:
	IMAGE_SERVICE=$(IMAGE_SERVICE) IMAGE_OUTBOX=$(IMAGE_OUTBOX) docker compose -v -f deploy/docker/dev-compose.yaml down