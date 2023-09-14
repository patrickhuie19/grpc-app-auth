# Github workflows

## Traces

`tracing.yml`:
This gh action will check if PRs are marked with the `enable tracing` label. If so, on PR creation, sycronization, or labeling, then a sample CI workflow is kicked off that generates traces.

Ephemeral OTEL collector, Grafana Tempo, and Grafana UI containers are deployed and kept up while the workflow executes (for 500 seconds). Ngrok tunnels are used to expose the Grafana UI container outside the github runner. 

To view such traces, fork this repository and configure the following secrets:
    NGROK_AUTH_TOKEN
    NGROK_USER
    NGROK_PASSWORD

The tracing workflow will log the exposed endpoint, which you can paste into any browser. Visit `/explore` and select traces to view traces.

`tracing-ssh.yml`:
This is a variant of `tracing.yml` that uses ssh reverse port forwarding and vanilla port forwarding instead of ngrok:

```mermaid
graph TD
  subgraph "GitHub Runner"
    A([Grafana UI<br>Container]) -->|Exposes Port 3000| B([Localhost<br>Port 3000])
  end
  subgraph "Remote Server"
    C([Port 3001])
  end
  subgraph "Authenticated User's Machine"
    E[View Grafana UI] -->|Browser Access| D([Localhost<br>Port 8000])
  end
  B .->|SSH Reverse Tunnel<br>Port 3000 -> 3001| C
  D -->|SSH Forward Tunnel<br>Port 3001 -> 8000| C

```
