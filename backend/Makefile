# ===================== Project Vars =====================
APP            := template-backend-go
PKG            := ./...
SERVER_MAIN    := ./cmd/server/main.go
BIN_DIR        := bin
OUT            := $(BIN_DIR)/server

GO             ?= go

# ---- Git metadata (Windows & *nix compatible) ----
VERSION        := $(shell git describe --tags --always --dirty 2>NUL || git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT         := $(shell git rev-parse --short HEAD 2>NUL || git rev-parse --short HEAD 2>/dev/null || echo none)

# ---- Date (ISO8601 UTC) ----
ifeq ($(OS),Windows_NT)
  DATE        := $(shell powershell -NoProfile -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ss'")Z
  SHELL       := bash
  RM_FILE     := cmd /C del /Q /F
  RM_DIR      := cmd /C rmdir /S /Q
  MKDIR_P     := mkdir
  NULLDEV     := NUL
  CHECKTOOL   = powershell -NoProfile -Command "if (-not (Get-Command $(1) -ErrorAction SilentlyContinue)) { exit 1 }"
else
  DATE        := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
  RM_FILE     := rm -f
  RM_DIR      := rm -rf
  MKDIR_P     := mkdir -p
  NULLDEV     := /dev/null
  CHECKTOOL   = command -v $(1) >/dev/null 2>&1
endif

LDFLAGS        := -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

# ---- DB / Migrations (goose) ----
MIGRATIONS_DIR := migrations
DB_DSN         ?= # example: postgres://user:pass@localhost:5432/db?sslmode=disable

# ---- Go proxy ----
GOPROXY_URL ?= https://proxy.golang.org,direct
GOSUMDB_URL ?= sum.golang.org

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

# ===================== Defaults =====================
.DEFAULT_GOAL := help

# ===================== Meta =====================
.PHONY: help tidy deps fmt vet lint test test-norace test-race cover run dev build clean \
        docker-build docker-run docker-stop \
        migrate-up migrate-down migrate-reset migrate-status migrate-create \
        db-up db-down tools-install tools-update print-vars \
        test-unit test-integration test-e2e test-all test-db-up test-db-down \
        mocks audit

help: ## Show help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS=":.*?## "}; {printf "  %-22s %s\n", $$1, $$2}'

print-vars: ## Print variables (debug)
	@echo "APP=$(APP)"
	@echo "VERSION=$(VERSION)"
	@echo "COMMIT=$(COMMIT)"

# ===================== Go basics =====================
tidy: ## go mod tidy
	$(GO) mod tidy

deps: ## Download dependencies
	$(GO) mod download

fmt: ## Format code
	$(GO) fmt $(PKG)

vet: ## Go vet
	$(GO) vet $(PKG)

lint: ## Run golangci-lint
	@$(call CHECKTOOL,golangci-lint) || { echo "golangci-lint missing. Run 'make tools-install'."; exit 1; }
	golangci-lint run

audit: lint ## Run security audit and lint
	@echo ">> Running govulncheck..."
	@$(call CHECKTOOL,govulncheck) || $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# ---- Race detector ----
ifeq ($(OS),Windows_NT)
  RACE_FLAG :=
  SET_CGO   :=
else
  RACE_FLAG := -race
  SET_CGO   := CGO_ENABLED=1
endif

test: ## Run unit tests (race enabled if supported)
	$(SET_CGO) $(GO) test $(PKG) $(RACE_FLAG) -cover

test-unit: ## Run unit tests (short)
	$(SET_CGO) $(GO) test -v -short $(RACE_FLAG) $(PKG)

test-integration: ## Run integration tests
	$(SET_CGO) $(GO) test -v -tags=integration ./tests/integration/...

test-e2e: ## Run e2e tests
	$(SET_CGO) $(GO) test -v -tags=e2e ./tests/e2e/...

test-all: test-unit test-integration ## Run all tests

cover: ## Coverage report
	$(SET_CGO) $(GO) test $(PKG) $(RACE_FLAG) -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

run: ## Run directly
	$(GO) run $(SERVER_MAIN)

dev: ## Hot reload (air)
	@$(call CHECKTOOL,air) || { echo "air missing. Run 'make tools-install'"; exit 1; }
	air

build: ## Build binary
	@$(MKDIR_P) $(BIN_DIR)
	$(GO) build -ldflags "$(LDFLAGS) -s -w" -o $(OUT) $(SERVER_MAIN)
	@echo Built: $(OUT)

clean: ## Clean build artifacts
	-@$(RM_DIR) $(BIN_DIR) 2>$(NULLDEV)
	-@$(RM_FILE) coverage.out 2>$(NULLDEV)
	-@$(RM_FILE) coverage.html 2>$(NULLDEV)

# ===================== Code Gen / Mocks =====================
mocks: ## Generate/Update mocks
	@$(call CHECKTOOL,mockery) || { echo "mockery missing. Run 'make tools-install'"; exit 1; }
	mockery --all --dir internal/repository --output tests/mocks --outpkg mocks

# ===================== Docker =====================
docker-build: ## Build docker image
	docker build -t $(APP):$(VERSION) .

docker-run: ## Run docker container
	docker run --rm -p 8080:8080 --env-file .env $(APP):$(VERSION)

# ===================== Migrations =====================
GOBIN_PATH := $(shell go env GOBIN)
ifeq ($(GOBIN_PATH),)
  GOBIN_PATH := $(shell go env GOPATH)/bin
endif
GOOSE := "$(GOBIN_PATH)/goose"

migrate-up: ## goose up
	@test -n "$(DB_DSN)" || (echo "DB_DSN required"; exit 1)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down: ## goose down
	@test -n "$(DB_DSN)" || (echo "DB_DSN required"; exit 1)
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-create: ## Create migration: make migrate-create name=foo
	@test -n "$(name)" || (echo "name required"; exit 1)
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

# ===================== Tools =====================
tools-install: ## Install dev tools
	@echo ">> Installing dev tools..."
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/pressly/goose/v3/cmd/goose@latest
	$(GO) install github.com/air-verse/air@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/vektra/mockery/v2@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest

# ===================== Swagger =====================
swag-install:
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

swagger: swag-install ## Generate swagger docs
	swag init -g cmd/server/main.go -o docs --parseInternal --parseDependency