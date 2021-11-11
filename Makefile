## build: build the application and place the built app in the bin folder
build:
	go build -o bin/larry ./cmd/larry/.

## start: start container
start:
	docker-compose up -d

## run-dev: runs the application in dev mode
run-dev:
	go run ./cmd/larry/. -t golang -x 1 --safe-mode

## test: runs tests
test:
	go test -v ./... --cover

## compile: compiles the application for multiple environments and place the output executables under the bin folder
compile:
	# 64-Bit
	# FreeBDS
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/larry-freebsd-64 ./cmd/larry/.
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -o ./bin/larry-macos-64 ./cmd/larry/.
	# Linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/larry-linux-64 ./cmd/larry/.
	# Windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/larry-windows-64 ./cmd/larry/.

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'