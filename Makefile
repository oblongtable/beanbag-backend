DEV_COMPOSE_FILE=docker-compose-dev.yml
DEBUG_COMPOSE_FILE=docker-compose-debug.yml
TEST_COMPOSE_FILE=docker-compose-test.yml

# Get the GOPATH/bin directory
GOPATH_BIN := $(shell go env GOPATH)/bin

# Add GOPATH/bin to the PATH for this Makefile
export PATH := $(GOPATH_BIN):$(PATH)

# Check if sqlc is installed
ifeq (, $(shell which sqlc))
    $(error sqlc is not installed. Please run 'go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest')
	$(error This may be due to use of sudo. Please just run make, or if you absolutely need 'sudo', use 'sudo -E make ...' to preserve your environment.)
endif

ifeq (, $(shell which swag))
	$(error swag could not be found. Please install it using: go install github.com/swaggo/swag/cmd/swag@latest')
endif
### DOCKER COMPOSE COMMANDS

.PHONY: compose-build
compose-build:
	sqlc generate
	swag init --pd
	sudo docker compose -f $(DEV_COMPOSE_FILE) build

.PHONY: compose-up
compose-up:
	sqlc generate
	swag init --pd
	sudo docker compose -f $(DEV_COMPOSE_FILE) up

.PHONY: compose-up-build
compose-up-build:
	sqlc generate
	swag init --pd
	sudo docker compose -f $(DEV_COMPOSE_FILE) up --build

.PHONY: compose-up-debug-build
compose-up-debug-build:
	sqlc generate
	swag init --pd
	sudo docker compose -f $(DEV_COMPOSE_FILE) -f $(DEBUG_COMPOSE_FILE) up --build

.PHONY: compose-down
compose-down:
	sudo docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: compose-down-wipe
compose-down-wipe:
	sudo docker compose -f $(DEV_COMPOSE_FILE) down -v

DOCKERCONTEXT_DIR:=./
DOCKERFILE_DIR:=./

.PHONY: run-tests
run-tests:
	sqlc generate
	swag init --pd
	sudo docker compose -f $(DEV_COMPOSE_FILE) -f $(TEST_COMPOSE_FILE) run --rm --build beanbag-backend
