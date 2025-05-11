#!/bin/bash

# Copy development env files
cp backend/config/.env.dev backend/config/.env
cp frontend/.env.dev frontend/.env

# Run docker-compose with development configuration
docker-compose up