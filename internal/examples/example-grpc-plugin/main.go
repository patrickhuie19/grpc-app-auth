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

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command("sh", "-c", os.Getenv("ADD_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.AddPluginName)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
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
