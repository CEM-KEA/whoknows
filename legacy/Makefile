.PHONY: init run test install clean

# Define environment variables
PYTHONPATH := ./src/backend
DATABASE_PATH := ./db/whoknows.db
TEST_DATABASE_PATH := ./db/test_whoknows.db
VENV := .venv

# Initialize the database
init:
	@echo "Initializing database..."
	PYTHONPATH=$(PYTHONPATH) $(VENV)/bin/python3 -c "from database import init_db; init_db()"
	@echo "Database initialized."

# Install dependencies and set up virtual environment
install:
	@echo "Setting up virtual environment and installing dependencies..."
	python3 -m venv $(VENV)
	$(VENV)/bin/pip install -r ./requirements.txt
	@echo "Dependencies installed."

# Run the Flask app
run:
	@echo "Checking if database exists..."
	PYTHONPATH=$(PYTHONPATH) $(VENV)/bin/python3 -c "from database import check_db_exists; check_db_exists()" || make init
	@echo "Running the Flask application..."
	PYTHONPATH=$(PYTHONPATH) $(VENV)/bin/python3 ./src/backend/app.py

# Run tests using unittest
test:
	@echo "Running tests..."
	PYTHONPATH=$(PYTHONPATH) DATABASE_PATH=$(TEST_DATABASE_PATH) $(VENV)/bin/python3 -m unittest discover -s tests -v
	@echo "Tests completed."

# Clean up virtual environment and any unnecessary files
clean:
	@echo "Cleaning up..."
	rm -rf $(VENV)
	rm -rf __pycache__
	rm -rf legacy/src/backend/__pycache__
	rm -rf tests/__pycache__
	rm -rf $(TEST_DATABASE_PATH)
	@echo "Clean complete."