# Distributed Tracing Example: OTEL + Grafana

This example deploys the grpc-app-auth server, an OTEL collector, grafana Tempo, and the grafana UI to containers on the same network.

A localhost client can send gRPC calls to the server. The gRPC server is instrumented with open telemetry traces, which are sent to the OTEL collector and forwarded to the Tempo backend. The grafana UI can then read the trace data from the Tempo backend. 

To get started:

1. Start up the stack
    ```bash
    docker-compose up --build
    ```

    You should see 4 containers started.

2. Send client calls

    Open a new terminal and run:

    ```bash
    cd ../example-grpc-client-server/client
    go run .
    ```

    To observe traces from the plugin example, navigate to `example-grpc-plugin`.

    Then follow the build instructions in [plugin README](../example-grpc-plugin/README.md)

    Run the example:
    ```bash
    ./add-service 1 2
    ```

3. Navigate to [Grafana](http://localhost:3000/explore) select the Tempo data source and use the "Search" tab to find traces.

    You should see traces from the server in `example-grpc-client-server` with the service name `grpc-app-auth-server` and from the plugin with the service name `add-plugin`

4. Teardown

    To stop the containers navigate back to the terminal where you ran `docker-compose up --build` and hit ctrl + c, or navigate to Docker Desktop and stop and delete the compose stack called `example-otlp-agent-tempo-grafana`.