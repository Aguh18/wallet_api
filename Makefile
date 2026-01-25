ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

BASE_STACK = docker compose -f docker-compose.yml
INTEGRATION_TEST_STACK = $(BASE_STACK) -f docker-compose-integration-test.yml
ALL_STACK = $(INTEGRATION_TEST_STACK)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
run: deps ### run application
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

dev: ### run application with hot reload using air
	air
.PHONY: dev

##@ Docker
compose-up: ### Run docker compose (database only) in background
	$(BASE_STACK) up --build -d db
.PHONY: compose-up

compose-up-all: ### Run docker compose (database + app)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-up-integration-test: ### Run docker compose with integration test
	$(INTEGRATION_TEST_STACK) up --build --abort-on-container-exit --exit-code-from integration-test
.PHONY: compose-up-integration-test

compose-down: ### Down docker compose
	$(ALL_STACK) down --remove-orphans
.PHONY: compose-down

nuke: ### Nuke docker - remove all containers, volumes, and networks
	docker compose down -v --remove-orphans
	docker volume prune -f
	docker system prune -f
	@echo "✅ Docker nuked! All containers, volumes, and networks removed."
.PHONY: nuke

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

##@ Dependencies
deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

deps-audit: ### check dependencies vulnerabilities
	govulncheck ./...
.PHONY: deps-audit

bin-deps: ### install development tools
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
.PHONY: bin-deps

##@ Database
migrate-create:  ### create new migration (usage: make migrate-create NAME=name)
	migrate create -ext sql -dir migrations $(NAME)
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

seed: ### run database seeder
	CGO_ENABLED=0 go run -tags migrate ./cmd/seed
.PHONY: seed

migrate-down: ### migration down (1 step)
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' down 1
.PHONY: migrate-down

migrate-down-all: ### migration down (all)
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' drop
.PHONY: migrate-down-all

##@ Testing
test: ### run test
	go test -v -coverprofile=coverage.txt -covermode=atomic ./internal/... ./pkg/...

test-race: ### run test with race detection
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./internal/... ./pkg/...
.PHONY: test test-race

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

##@ Code Quality
format: ### Run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

pre-commit: format linter-golangci test ### run pre-commit checks
.PHONY: pre-commit
