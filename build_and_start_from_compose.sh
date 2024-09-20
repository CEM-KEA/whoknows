#!/bin/bash

# Check if the external volume exists
if ! docker volume ls --format '{{.Name}}' | grep -q "^sqlite_data$"; then
  echo "Creating external volume 'sqlite_data'..."
  docker volume create sqlite_data
else
  echo "Volume 'sqlite_data' already exists."
fi

# Build and start the application with Docker Compose
docker compose up --build -d