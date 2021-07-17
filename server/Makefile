# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
.PHONY: run
run:
	 go run gotodo.rasc.ch/cmd/web

.PHONY: build/goose
build/goose:
	@go build -o migrate gotodo.rasc.ch/cmd/migrate

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

## db/migration/status: show database migration status
.PHONY: db/migration/status
db/migration/status: build/goose
	./migrate status
	
## db/codegen: generate sqlboiler code
.PHONY: db/codegen
db/codegen:
	@docker build -t sqlboilercodegen -f sqlboiler/Dockerfile .
	@docker run -v $(shell pwd):/src sqlboilercodegen

.PHONY: db/run-sqlboiler
db/run-sqlboiler:
	cd sqlboiler && sqlboiler mysql

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## lint: run golangci-lint
.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.41.1 golangci-lint run ./...

## upgrade-libraries: upgrade all dependant libraries
.PHONY: upgrade-libraries
upgrade-libraries:
	@go get -u ./...
	@go fmt ./...
	@go mod tidy
	@go mod verify

# ==================================================================================== #
# BUILD
# ==================================================================================== #
.PHONY: build/docker
build/docker:
	@docker build -t gobuilder -f docker/build.Dockerfile .

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.appBuildTime=${current_time} -X main.appVersion=${git_description}'

.PHONY: build/web
build/web: build/docker
	@echo 'Building cmd/web for ${goos}...'
	@docker run -v $(shell pwd):/src gobuilder env GOOS=${goos} GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/${goos}/gotodo${extension} gotodo.rasc.ch/cmd/web

## build/web/linux: build the cmd/web application for linux
.PHONY: build/web/linux
build/web/linux: goos=linux
build/web/linux: extension=
build/web/linux: build/web

## build/web/windows: build the cmd/web application for windows
.PHONY: build/web/windows
build/web/windows: goos=windows
build/web/windows: extension=.exe
build/web/windows: build/web
