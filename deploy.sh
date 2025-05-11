#!/bin/bash

# Copy production env files
cp backend/config/.env.prod backend/config/.env
cp frontend/.env.prod frontend/.env

# Run docker-compose with production configuration
docker-compose up -d