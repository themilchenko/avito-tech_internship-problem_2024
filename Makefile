MAIN_PATH = ./cmd/main.go

.PHONY: local_build
local_build: ## Build locally
	go build -o bin/ ${MAIN_PATH}

.PHONY: docker_build
docker_build: ## Build project in docker compose
	docker container prune -f
	docker-compose up -d

.PHONY: lint
lint: ## Make linters
	@golangci-lint run -c configs/.golangci.yaml

.PHONY: e2e
e2e: ## Make end-to-end tests
	ADMIN_USERNAME=milchenko ADMIN_PASSWORD=qwerty USER_USERNAME=user USER_PASSWORD=qwerty pytest e2e/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := local_build
