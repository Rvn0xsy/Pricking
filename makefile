GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: build

build:
	@echo "Building Pricking...."
	$(GOBUILD) -v -ldflags="-extldflags=-static" -o "pricking" cmd/pricking/pricking.go

test:
	@./pricking -url https://payloads.online -config ./config/config.yaml
clean:
	@rm -rf ./pricking
help:
	@echo make build
	@echo make test
	@echo make test
