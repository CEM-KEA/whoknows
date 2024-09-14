# Makefile

# Variables
BACKEND_BUILD_SCRIPT = build.bash
FRONTEND_BUILD_CMD = cd ./frontend/ && npm run build
FRONTEND_DEV_CMD = cd ./frontend/ && npm run dev
BACKEND_TEST_CMD = cd ./backend/ && go test ./...
FRONTEND_TEST_CMD = cd ./frontend/ && npx playwright test

# Target to build only the backend
backend-build:
	@echo "Building backend..."
	@./$(BACKEND_BUILD_SCRIPT)

# Target to run backend tests
backend-test:
	@echo "Running backend tests..."
	@$(BACKEND_TEST_CMD)

# Target to build only the frontend
frontend-build:
	@echo "Building frontend..."
	@$(FRONTEND_BUILD_CMD)

frontend-test:
	@echo "Running frontend tests..."
	@$(FRONTEND_TEST_CMD)

# Target to build and run backend, then build and run frontend
build:
	@echo "Building and running backend and frontend..."
	@./$(BACKEND_BUILD_SCRIPT) &
	@$(FRONTEND_BUILD_CMD)

# Target for development mode - build and run backend, and run frontend in dev mode
dev:
	@echo "Starting development mode..."
	@./$(BACKEND_BUILD_SCRIPT) &
	@$(FRONTEND_DEV_CMD)
