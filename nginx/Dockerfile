# Use the official Node.js image as the base for the build stage
FROM node:24-slim AS build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY ../frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the frontend files
COPY ../frontend .

# Build the frontend application
RUN npm run build

# Use the official Nginx image as the base for the production image
FROM nginx:1.28.0-alpine

# Install Certbot
RUN apk add --no-cache certbot certbot-nginx

# Copy the built frontend files from the build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Remove the default Nginx configuration file
RUN rm /etc/nginx/conf.d/default.conf

# Copy the custom Nginx configuration file
COPY nginx/nginx.conf /etc/nginx/nginx.conf

# Expose ports 80 and 443
EXPOSE 80 443

# Command to start Nginx
CMD ["nginx", "-g", "daemon off;"]