# Makefile

# Variables
BACKEND_BUILD_SCRIPT = build.bash
FRONTEND_DEV_CMD = cd ./frontend/ && npm run dev
BACKEND_TEST_CMD = cd ./backend/ && go test ./...
FRONTEND_TEST_CMD = cd ./frontend/ && npx playwright test
DOCKER_COMPOSE_DOWN_CMD = docker compose down
DOCKER_MAKE_EXTERNAL_VOLUME_SQLITE_DATA = docker volume create --name=sqlite_data
DOCKER_COMPOSE_BUILD_CMD = docker compose build --no-cache
DOCKER_COMPOSE_CMD = docker compose up --force-recreate

backend-dev:
	@echo "Starting backend in development mode..."
	@./$(BACKEND_BUILD_SCRIPT)

frontend-dev:
	@echo "Starting frontend in development mode..."
	@$(FRONTEND_DEV_CMD)

backend-test:
	@echo "Running backend tests..."
	@$(BACKEND_TEST_CMD)

frontend-test: compose-detach
	@echo "Running frontend tests..."
	@$(FRONTEND_TEST_CMD)
	@$(DOCKER_COMPOSE_DOWN_CMD)

compose-down:
	@echo "Stopping backend and frontend..."
	@$(DOCKER_COMPOSE_DOWN_CMD)

build-images: compose-down
	@echo "Building backend and frontend images..."
	@$(DOCKER_MAKE_EXTERNAL_VOLUME_SQLITE_DATA)
	@$(DOCKER_COMPOSE_BUILD_CMD)

compose: build-images
	@echo "Composing backend and frontend..."
	@$(DOCKER_COMPOSE_CMD)

compose-detach: build-images
	@echo "Composing backend and frontend in detached mode..."
	@$(DOCKER_COMPOSE_CMD) -d

dev:
	@echo "Starting development mode..."
	@$(MAKE) backend-dev &
	@$(MAKE) frontend-dev

test:backend-test frontend-test
	
