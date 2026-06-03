.PHONY: swag

# Container name (from docker-compose.yml)
CONTAINER_NAME_API := api
CONTAINER_NAME_ADMIN := admin

# Check if container is running
CONTAINER_RUNNING := $(shell docker ps --format '{{.Names}}' 2>/dev/null | grep -E '^$(CONTAINER_NAME_API)$$')

# Determine execution command
ifdef CONTAINER_RUNNING
	EXEC_CMD = docker exec $(CONTAINER_NAME_API)
	EXEC_CMD_INTERACTIVE = docker exec -i $(CONTAINER_NAME_API)
	ENV_MSG = 📦 Running in Docker container
else
	EXEC_CMD = 
	EXEC_CMD_INTERACTIVE = 
	ENV_MSG = 💻 Running on host (Docker not available)
endif

## swag: Generate Swagger documentation
swag:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) swag init -g ./cmd/server/main.go -o ./api/docs
	@echo "✅ Swagger docs generated"
else
	@if command -v swag >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		swag init -g ./cmd/server/main.go -o ./api/docs; \
		echo "✅ Swagger docs generated"; \
	else \
		echo "❌ Error: Docker container not running and swag not installed"; \
		echo "Please run: make up"; \
		echo "Or install: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi
endif
