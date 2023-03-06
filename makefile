APP=ast-app-agent

.PHONY: help all build windows linux darwin

help:
	@echo "usage: make <option>"
	@echo "options and effects:"
	@echo "    help   : Show help"
	@echo "    all    : Build multiple binary of this project"
	@echo "    build  : Build the binary of this project for current platform"
	@echo "    windows: Build the windows binary of this project"
	@echo "    linux  : Build the linux binary of this project"
	@echo "    darwin : Build the darwin binary of this project"

all:build windows linux darwin
build:
	@echo "start build ${APP} ..."
	go build -o build/ci/bin/${APP}

windows:
	CGO_ENABLED=0 GOOS=windows go build -o build/ci/bin/${APP}-windows

linux:
	CGO_ENABLED=0 GOOS=linux go build -o build/ci/bin/${APP}-linux

darwin:
	CGO_ENABLED=0 GOOS=darwin go build -o bin/${APP}-darwin
