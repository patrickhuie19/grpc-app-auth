package main

import (
	"log"
	"grpc-app-auth/internal/examples/example-grpc-plugin/shared"
	"grpc-app-auth/server"
	"os"

	"github.com/hashicorp/go-plugin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Implementation of AddInterface that implements the Add business logic
type AddInterface struct{}

func (AddInterface) Add(a float64, b float64) (float64, error) {
	return a + b, nil
}

func main() {
	enableTelemetry := os.Getenv("ENABLE_TELEMETRY")
	telemetryTarget := os.Getenv("TELEMETRY_TARGET")

	if enableTelemetry == "true" {
		server := &server.Server{}
		err := server.SetupOpenTelemetry(telemetryTarget, "plugin-add"); if err != nil {
			log.Printf("Plugin Error: %v", err.Error())
			os.Exit(1)
		}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.AddPluginName: &shared.AddGRPCPlugin{Impl: &AddInterface{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
            opts = append(opts, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
            opts = append(opts, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
            return grpc.NewServer(opts...)
		},
	})
}
