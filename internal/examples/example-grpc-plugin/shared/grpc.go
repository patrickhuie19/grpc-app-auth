package shared

import (
	"context"

	"grpc-app-auth/services"
)

// GRPCClient is an implementation of AddInterface that talks over RPC.
type GRPCClient struct{ client services.AddClient }

func (m *GRPCClient) Add(a float64, b float64) (float64, error) {
	result, err := m.client.Add(context.Background(), &services.AddRequest{
		A: a,
		B: b,
	})
	if err != nil {
		return float64(0), err
	}
	return result.Result, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl AddInterface

	services.UnimplementedAddServer
}

func (m *GRPCServer) Add(
	ctx context.Context,
	req *services.AddRequest) (*services.AddReply, error) {
	v, err := m.Impl.Add(req.A, req.B)
	return &services.AddReply{Result: v}, err
}
