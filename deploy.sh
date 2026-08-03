#!/bin/bash

echo "Stopping containers"
docker compose down

echo "Fetching latest changes from git"
git pull

echo "Pulling any changes from Docker images"
docker compose --profile prod pull

echo "Rebuilding contianers"
docker compose up -d --build