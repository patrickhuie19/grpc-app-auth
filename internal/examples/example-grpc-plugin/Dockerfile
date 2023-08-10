FROM golang:latest

# Setting up the working directory inside the container.
WORKDIR /app

# Copy all files from the current directory to the container's working directory.
COPY . .

# Switch to the directory containing the main CLI.
WORKDIR /app/internal/examples/example-grpc-plugin

# Build the main CLI and the plugin.
RUN go build -o add-service
RUN go build -o add-service-grpc ./add-plugin

# Set the environment variable in the container.
ENV ADD_PLUGIN="./add-service-grpc"

# When the container starts, run the add-service binary.
ENTRYPOINT ["./add-service", "1", "2"]
