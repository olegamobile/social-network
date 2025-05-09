#!/bin/bash

# Copy development env files
cp backend/.env.dev backend/.env
cp frontend/.env.dev frontend/.env

# Run docker-compose with development configuration
docker-compose up