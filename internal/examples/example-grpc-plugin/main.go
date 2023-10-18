package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"grpc-app-auth/internal/examples/example-grpc-plugin/shared"
	"github.com/hashicorp/go-plugin"
)

func main() {
	enableTelemetry := os.Getenv("ENABLE_TELEMETRY")
	telemetryTarget := os.Getenv("TELEMETRY_TARGET")

	// We're a host. Start by launching the plugin process.
	cmd := exec.Command("sh", "-c", os.Getenv("ADD_PLUGIN"))
	cmd.Env = []string{"ENABLE_TELEMETRY=" + enableTelemetry, "TELEMETRY_TARGET=" + telemetryTarget}
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via GRPC
	grpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := grpcClient.Dispense(shared.AddPluginName)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	// We should have an add service now! This feels like a normal interface
	// implementation but is in fact over a GRPC connection.
	addService := raw.(shared.AddInterface)
	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Malformed input: %v", err.Error())
		os.Exit(1)
	}

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Printf("Malformed input: %v", err.Error())
		os.Exit(1)
	}
	result, err := addService.Add(a, b)
	if err != nil {
		fmt.Printf("Add service failed: %v", err.Error())
		os.Exit(1)
	}

	fmt.Println(result)

	os.Exit(0)
}
