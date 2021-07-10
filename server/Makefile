run:
	 go run gotodo.rasc.ch/cmd/web

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.41.1 golangci-lint run ./...

build-prod:
	go build -ldflags='-s' -o gotodo cmd/web

upgrade-libraries:
	@go get -u ./...
	@go fmt ./...
	@go mod tidy
	@go mod verify
