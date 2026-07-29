#!/bin/bash

echo "Fetching latest changes from git"
git pull

echo "Stopping containers"
docker compose down

echo "Rebuilding contianers"
docker compose -d --build