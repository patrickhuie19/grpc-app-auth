FROM golang:latest

WORKDIR /app

COPY . .

WORKDIR /app/internal/examples/example-grpc-client-server/server/

RUN go mod download
RUN go clean
RUN go build -o server server_main.go
RUN chmod +x server

# CMD ["/bin/sh"]
CMD ["./server"]
