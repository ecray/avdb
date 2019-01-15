# Command parameters
GOCMD=go
DEP=dep
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=release/avdb
BINARY_LINUX=$(BINARY_NAME)_linux_amd64
REGISTRY=avdb

all: test build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_LINUX)
		docker-compose down
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
deps:
		$(DEP) init
		$(DEP) ensure

# Release container
#dev-env: build-linux docker-build docker-compose
dev-env: build-linux docker-compose

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
docker-compose:
	docker-compose build
	docker-compose up

# Tag container and publish to registry
publish-release: build-linux image-build tag publish

image-build:
	docker build -t avdb . --no-cache
tag:
	docker tag avdb:latest ${REGISTRY}
publish:
	docker push ${REGISTRY}
