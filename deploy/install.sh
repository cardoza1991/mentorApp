#!/usr/bin/env bash

set -e

# Function to print messages
function echo_info {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

function echo_error {
    echo -e "\033[1;31m[ERROR]\033[0m $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null
then
    echo_error "Docker is not installed. Please install Docker and try again."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null
then
    echo_error "Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

# Validation: Ensure running within Athena0
# Example validation: Check if the IP address 10.20.0.1 is assigned to any network interface
if ! ip addr show | grep -q "10.20.0.1"; then
    echo_error "This installer must be run within the Athena0 environment (IP: 10.20.0.1 is not found)."
    exit 1
fi

echo_info "Validation passed. Proceeding with deployment."

# Navigate to the directory where the script is located
cd "$(dirname "$0")"

# Deploy the Docker Compose stack
echo_info "Deploying Docker Compose stack..."
docker-compose up -d

# Wait for services to be healthy
echo_info "Waiting for services to become healthy..."
for i in {1..10}; do
    HEALTH_STATUS=$(docker-compose ps -q postgres | xargs docker inspect -f '{{.State.Health.Status}}')
    if [[ "$HEALTH_STATUS" == "healthy" ]]; then
        echo_info "Postgres is healthy."
        break
    else
        echo_info "Waiting for Postgres to be healthy... ($i/10)"
        sleep 5
    fi

    if [[ "$i" -eq 10 ]]; then
        echo_error "Postgres failed to become healthy in time."
        exit 1
    fi
done

# Open the application in the default browser
echo_info "Opening the application in your default browser..."
xdg-open http://10.20.0.1:8081 || echo_error "Failed to open browser. Please navigate to http://10.20.0.1:8081 manually."

echo_info "Deployment completed successfully."
