FROM node:24-slim

WORKDIR /app

# Copy only the files needed for install/build
COPY package*.json ./

# Install dependencies before copying the whole context
RUN npm install

# Copy rest of the files after installing deps
COPY . .

# Default API URL for development; can be overridden
ENV VITE_API_URL=http://localhost:8080

# Build the app
RUN npm run build

# Install `serve` to serve static files
RUN npm install -g serve

EXPOSE 5173

CMD ["serve", "-s", "dist", "-l", "5173"]
