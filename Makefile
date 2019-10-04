GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
GOCLEAN=$(GOCMD) clean
GOLINT=golint
GOBINARY=recipey_api
DOCKER_IMAGE_NAME=recipey_api
MIGRATE_CMD=migrate
MIGRATION_DIR=cmd/db_migration/migrations

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
upbuild:
	docker-compose up --build
down:
	docker-compose down
restart:
	docker-compose restart
restart-api:
	$(DC) restart api
bash-api:
	docker-compose exec api bash
bash-db:
	docker-compose exec db bash
migrate-create:
ifdef name
	$(MIGRATE_CMD) create -ext sql -dir $(MIGRATION_DIR) $(name)
else
	@echo "name=<name_of_your_migration> required as an argument"
endif
migrate-up:
	$(MIGRATE_CMD) -path $(MIGRATION_DIR) -database $(APP_DB_URL) up
migrate-down:
	$(MIGRATE_CMD) -path $(MIGRATION_DIR) -database $(APP_DB_URL) down
migrate-force:
	$(MIGRATE_CMD) -path $(MIGRATION_DIR) -database $(APP_DB_URL) force

.PHONY: server
