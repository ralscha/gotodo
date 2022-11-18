# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## lint: run golangci-lint
.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v

## upgrade-libraries: upgrade all dependant libraries
.PHONY: upgrade-libraries
upgrade-libraries:
	@go get -u ./...
	@go fmt ./...
	@go mod tidy
	@go mod verify


## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go mod verify


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cmd/web application
.PHONY: build
build:
	go mod verify
	go build -ldflags='-s' -o=./bin/web ./cmd/web
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/web ./cmd/web

## run: run the cmd/web application
.PHONY: run
run: tidy build
	./bin/web


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

.PHONY: build/goose
build/goose:
	go build -o migrate gotodo.rasc.ch/cmd/migrate

## db/migration/new/sql name=$1: create a new sql database migration
.PHONY: db/migration/new/sql
db/migration/new/sql: build/goose
	./migrate create ${name} sql

## db/migration/new/go name=$1: create a new go database migration
.PHONY: db/migration/new/go
db/migration/new/go: build/goose
	./migrate create ${name} go

## db/migration/up: apply all up database migrations
.PHONY: db/migration/up
db/migration/up: build/goose
	./migrate up

## db/migration/reset: revert all database migrations
.PHONY: db/migration/reset
db/migration/reset: build/goose
	./migrate reset

## db/migration/status: show database migration status
.PHONY: db/migration/status
db/migration/status: build/goose
	./migrate status

## db/codegen: generate sqlboiler code
.PHONY: db/codegen
db/codegen:
	@docker build -t sqlboilercodegen -f sqlboiler/Dockerfile .
	@docker run -v $(shell pwd):/src sqlboilercodegen
	go run gotodo.rasc.ch/cmd/tsgen

.PHONY: db/run-sqlboiler
db/run-sqlboiler:
	cd sqlboiler && sqlboiler mysql
