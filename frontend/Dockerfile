FROM node:23-slim

WORKDIR /app
COPY . .

# Setting default env variable
ENV VITE_API_URL=http://localhost:8080

# Install dependencies
RUN npm install

# Build for production
RUN npm run build

RUN npm install -g serve
EXPOSE 5173

# Start the static file server
CMD ["serve", "-s", "dist", "-l", "5173"]