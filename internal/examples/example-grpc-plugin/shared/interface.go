// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"
	"net/rpc"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"grpc-app-auth/services"
)

const AddPluginName = "add_grpc"

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	AddPluginName: &AddGRPCPlugin{},
}

// AddInterface is the interface that we're exposing as a plugin.
type AddInterface interface {
	Add(a float64, b float64) (float64, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type AddPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl AddInterface
}

func (p *AddPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*AddPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type AddGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl AddInterface
}

func (p *AddGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	services.RegisterAddServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *AddGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: services.NewAddClient(c)}, nil
}
