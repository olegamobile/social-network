# --- Builder stage ---
FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o server .

# --- Final image ---
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
COPY --from=builder /app/data ./data
COPY --from=builder /app/migrations ./migrations

ENV PORT=8080
ENV FRONTEND_URL=http://localhost:5173

EXPOSE 8080

CMD ["./server"]
