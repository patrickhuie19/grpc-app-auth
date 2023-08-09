# Distributed Tracing Example: OTEL + Grafana

This example deploys the grpc-app-auth server, an OTEL collector, grafana Tempo, and the grafana UI to containers on the same network. A localhost client can then send gRPC calls to the server. The gRPC server is instrumented with open telemetry traces, which are sent to the OTEL collector and forwarded to the Tempo backend. The grafana UI can then read the trace data from the Tempo backend. 

1. Start up the stack
```bash
docker-compose up --build
```

You should see 4 containers started.

2. Navigate to [Grafana](http://localhost:3000/explore) select the Tempo data source and use the "Search" tab to find traces. 

3. To stop the setup hit ctrl + c, or navigate to Docker Desktop and stop and delete the compose stack called `example-otlp-agent-tempo-grafana`


