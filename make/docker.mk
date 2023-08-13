##@ Docker Development

DOCKER_COMPOSE := docker-compose

# ==============================================================================
# Public Commands:
# ==============================================================================

.PHONY: docker-compose
docker-compose: ## Run containers in deploy/dockercompose/docker-compose.yaml
docker-compose: docker.compose

.PHONY: docker-compose-down
docker-compose-down: ## Stop and remove containers in deploy/dockercompose/docker-compose.yaml
docker-compose-down: docker.compose.down

# ==============================================================================
# Private Commands:
# ==============================================================================

.PHONY: docker.compose
docker.compose:
	$(DOCKER_COMPOSE) -f deploy/dockercompose/docker-compose.yaml up -d

.PHONY: docker.compose.down
docker.compose.down:
	$(DOCKER_COMPOSE) -f deploy/dockercompose/docker-compose.yaml down