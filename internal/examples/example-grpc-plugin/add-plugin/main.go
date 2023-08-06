package main

import (
	"github.com/hashicorp/go-plugin"
	"grpc-app-auth/internal/examples/example-grpc-plugin/shared"
)

// Implementation of AddInterface that implements the Add business logic
type AddInterface struct{}

func (AddInterface) Add(a float64, b float64) (float64, error) {
	return a + b, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.AddPluginName: &shared.AddGRPCPlugin{Impl: &AddInterface{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
