.PHONY: dev

DEV_STATUS_FILE := .dev_status

dev:
	@if [ ! -f $(DEV_STATUS_FILE) ]; then \
		echo "Starting development environment..."; \
		docker compose -f compose.dev.yml up --build; \
		touch $(DEV_STATUS_FILE); \
	else \
		echo "Stopping development environment..."; \
		docker compose -f compose.dev.yml down; \
		rm $(DEV_STATUS_FILE); \
	fi