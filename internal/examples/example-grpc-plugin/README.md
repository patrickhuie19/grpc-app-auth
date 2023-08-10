# Plugin Add Service Example

Adapted heavily from https://github.com/hashicorp/go-plugin/tree/main/examples/grpc

## To build and run this example
```bash
# This builds the main CLI
$ go build -o add-service

# This builds the plugin written in Go
$ go build -o add-service-grpc ./add-plugin

# This tells the add-service binary to use the "add-service-grpc" binary
$ export ADD_PLUGIN="./add-service-grpc"

# Perform Add
$ ./add-service 1 2
```

## Performance
Add service over plugin is ~10x e2e lower latency than over the straight gRPC client server implemented in `example-grpc-client-server`