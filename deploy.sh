#!/bin/bash

# Check if .env exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create .env file from .env.template"
    exit 1
fi

# Load environment variables
source .env

# Pull latest docker-compose.yml
wget https://raw.githubusercontent.com/${GITHUB_USERNAME}/${REPO_NAME}/main/docker-compose.yml

# Start or update services
docker-compose pull
docker-compose up -d