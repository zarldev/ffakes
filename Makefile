.PHONY: build build-all build-mac build-linux build-windows install lint test generate help
default: help

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build: generate test lint
	go mod tidy && go build $(LDFLAGS) --tags prod -o dist/ffakes ffakes.go

build-all: build-mac build-linux build-windows

build-mac: build
	GOOS=darwin GOARCH=amd64 go mod tidy && go build $(LDFLAGS) --tags prod -o dist/ffakes ffakes.go

build-linux: build
	GOOS=linux GOARCH=amd64 go mod tidy && go build $(LDFLAGS) --tags prod -o dist/ffakes ffakes.go

build-windows: build
	GOOS=windows GOARCH=amd64 go mod tidy && go build $(LDFLAGS) --tags prod -o dist/ffakes.exe ffakes.go

releases: build-all
	@echo "Creating releases"
	@echo "Creating release for mac"
	tar -czvf dist/ffakes.mac.$(VERSION).tar.gz dist/ffakes
	@echo "Creating release for linux"
	tar -czvf dist/ffakes.linux.$(VERSION).tar.gz dist/ffakes
	@echo "Creating release for windows"
	zip -r dist/ffakes.windows.$(VERSION).zip dist/ffakes.exe

install:
	chmod +x bin/ffakes
	cp bin/ffakes /usr/local/go/bin/ffakes

lint:
	golangci-lint run ./...
	
test:
	go test -v ./...

generate:
	go generate ./...

clean:
	rm -rf dist
	go clean --testcache

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        	 Build the project for the current platform"
	@echo "  build-all    	 Build the project for all platforms"
	@echo "  build-mac    	 Build the project for mac"
	@echo "  build-linux  	 Build the project for linux"
	@echo "  build-win      Build the project for windows"
	@echo "  releases     	 Create releases for all platforms"
	@echo "  install      	 Install the project"
	@echo "  lint         	 Run linter"
	@echo "  test         	 Run tests"
	@echo "  clean        	 Clean the project"
	@echo "  generate     	 Generate code"
	@echo "  help         	 Show this help message"