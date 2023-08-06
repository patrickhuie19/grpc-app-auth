package shared

import (
	"net/rpc"
)

// RPCClient is an implementation of AddInterface that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) Add(a float64, b float64) (float64, error) {
	var resp float64
	err := m.client.Call("Plugin.Add", map[string]interface{}{
		"a": a,
		"b": b,
	}, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl AddInterface
}

func (m *RPCServer) Add(args map[string]interface{}, resp *interface{}) (float64, error) {
	return m.Impl.Add(args["a"].(float64), args["b"].(float64))
}
