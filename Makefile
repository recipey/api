GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
GOCLEAN=$(GOCMD) clean
GOLINT=golint
GOBINARY=recipey_api
DOCKER_IMAGE_NAME=recipey_api

all: clean get check docker-build build
check: lint vet
lint: ## Lint all files
	$(GOLIST) ./... | grep -v /vendor/ | xargs $(GOLINT) -set_exit_status=1
vet: ## Vet all files
	$(GOCMD) vet $(shell $(GOLIST) ./... | grep -v /vendor/)
get:
	$(GOLIST) ./... | grep -v /vendor/ | xargs $(GOGET)
build:
	$(GOBUILD) -o ./cmd/server/server -v ./cmd/server/
build-linux:
	GOOS=linux $(GOBUILD) -o ./cmd/server/server -v ./cmd/server/
server:
	./cmd/server/server
clean:
	$(GOLIST) ./... | xargs $(GOCLEAN)
docker-build:
	docker build . -f Dockerfile.api -t ${DOCKER_IMAGE_NAME}
up:
	docker-compose up
down:
	docker-compose down
restart:
	docker-compose restart

.PHONY: server
