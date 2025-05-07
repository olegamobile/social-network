FROM golang:1.24

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o server .

# Environment variables will be passed from docker-compose or docker run
ENV PORT=8080
ENV FRONTEND_URL=http://localhost:5173

EXPOSE 8080
CMD ["./server"]