#!/bin/bash

# Copy production env files
cp backend/.env.prod backend/.env
cp frontend/.env.prod frontend/.env

# Run docker-compose with production configuration
docker-compose up -d