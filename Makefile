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

# Runs the backend in development mode
backend-dev:
	@echo "Starting backend in development mode..."
	@./$(BACKEND_BUILD_SCRIPT)

# Runs the frontend in development mode
frontend-dev:
	@echo "Starting frontend in development mode..."
	@$(FRONTEND_DEV_CMD)

# Runs the backend tests
backend-test:
	@echo "Running backend tests..."
	@$(BACKEND_TEST_CMD)

# Runs docker compose to start the backend and frontend containers in detached mode
# Then runs the frontend/end-2-end tests
# Finally, runs docker compose down to stop the containers
frontend-test: compose-detach
	@echo "Running frontend tests..."
	@$(FRONTEND_TEST_CMD)
	@$(DOCKER_COMPOSE_DOWN_CMD)

# Stops the backend and frontend containers
compose-down:
	@echo "Stopping backend and frontend..."
	@$(DOCKER_COMPOSE_DOWN_CMD)

# Makes the external volume for the sqlite data required by backend image
# Builds the backend and frontend docker images
build-images: compose-down
	@echo "Building backend and frontend images..."
	@$(DOCKER_MAKE_EXTERNAL_VOLUME_SQLITE_DATA)
	@$(DOCKER_COMPOSE_BUILD_CMD)

# Builds the backend and frontend docker images
# Then runs docker compose to start the backend and frontend containers
compose: build-images
	@echo "Composing backend and frontend..."
	@$(DOCKER_COMPOSE_CMD)

# Builds the backend and frontend docker images
# Then runs docker compose to start the backend and frontend containers in detached mode
compose-detach: build-images
	@echo "Composing backend and frontend in detached mode..."
	@$(DOCKER_COMPOSE_CMD) -d

# Runs the backend and frontend in development mode
dev:
	@echo "Starting development mode..."
	@$(MAKE) backend-dev &
	@$(MAKE) frontend-dev

# Runs the backend and frontend tests
test:backend-test frontend-test
	
