# grpc-app-auth

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](LINK-TO-BUILD)
[![Go Version](https://img.shields.io/badge/go-%5E1.16-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/license/mit/)

## Introduction

`grpc-app-auth` is a demo project that showcases the capabilities of ed25519 public-key signatures, distributed tracing, and secure RPC calls. It allows users to run a server and client example to explore the technology and understand its implementation.

## Quick Installation

```bash
git clone https://github.com/patrickhuie19/grpc-app-auth.git
cd grpc-app-auth/internal/example
```

### Running the Server and Client

To give the code a try for the first time, you can run the server and client in two different terminals.

#### Server

Navigate to the server directory and run:

```
cd internal/example/server
go run .
```

#### Client

Similarly, for the client:

```bash
cd internal/example/client
go run .
```

## Allowlist Trust Store

If you want to understand how the allowlist trust store works, you can check out the integration tests in the intgtest directory and run the tests with:

```bash
cd internal/intgtest
go test ./... -v
```

## Distributed Tracing

We are experimenting with instrumenting this application with OpenTelemetry tracing. Check out branches other than main to see the ongoing work in this area.

## Contributions

Contributions and PRs are most welcome! Feel free to fork the repository, make your changes, and submit a pull request.


## Roadmap

### Signature-Based RPC Calls
Create an example RPC call that uses signatures in the gRPC request metadata, and not in the message body

## License

This project is licensed under the terms of the MIT license