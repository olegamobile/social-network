FROM golang:1.20

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o server .

RUN mkdir -p ./database
COPY ./database/init.sql ./database/

RUN sqlite3 ./database/app.db < ./database/init.sql

EXPOSE 8080
CMD ["./server"]
