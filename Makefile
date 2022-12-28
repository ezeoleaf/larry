## build: build the application (for current GOOS/GOARCH) and place the built app in the dist folder
build:
	goreleaser build --single-target --snapshot --rm-dist

## start: start container
start:
	docker-compose up -d

## run-dev: runs the application in dev mode
run-dev:
	go run ./cmd/larry/. -t golang -x 1 --safe-mode

## test: runs tests
test:
	go test -v ./... --cover

## compile: compiles the application for multiple environments and place the output under the dist folder
compile:
	goreleaser release --snapshot --rm-dist

lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
