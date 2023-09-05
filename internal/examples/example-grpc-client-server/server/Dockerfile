FROM golang:latest

WORKDIR /app

COPY . .

RUN ls -al

RUN go mod download

WORKDIR /app/internal/examples/example-grpc-client-server/server/

RUN go clean
RUN go build -o server server_main.go
RUN chmod +x server

CMD ["./server"]
