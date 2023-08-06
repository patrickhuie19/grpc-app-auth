# grpc-app-auth

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](LINK-TO-BUILD)
[![Go Version](https://img.shields.io/badge/go-%5E1.16-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/license/mit/)

## Introduction

`grpc-app-auth` is a demo project that showcases the capabilities of ed25519 public-key signatures, distributed tracing, and secure RPC calls.

## Quick Installation

```bash
git clone https://github.com/patrickhuie19/grpc-app-auth.git
cd grpc-app-auth/internal/example
```

## Examples
Checkout the examples directory for different ways of interacting with the implemented gRPC services

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

### Containerization

We're working on providing relevant docker templates and documentation to run the examples from docker containers. 

## License

This project is licensed under the terms of the MIT license