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

## Server and Client API

Check out the integration tests in the intgtest directory and run the tests to get a feel for the Client and Server API.

```bash
cd internal/intgtest
go test ./... -v
```

## Distributed Tracing

The gRPC services in this application can be instrumented OpenTelemetry tracing. Check out `example-otlp-agent-tempo-grafana`.

## Contributions

Contributions and PRs are most welcome! Feel free to fork the repository, make your changes, and submit a pull request.


## Roadmap

Nothing at the moment - feel free to suggest!