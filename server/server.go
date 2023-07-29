package server

import (
	"context"
	"crypto/ed25519"
	"log"
	"net"

	pb "grpc-app-auth/echo"

	"grpc-app-auth/internal/keystore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedEchoServiceServer
	trustedKeys keystore.KeyStore
	grpcServer *grpc.Server
}

func (s *Server) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoReply, error) {
	pubKey, err := s.trustedKeys.GetPublicKey(in.PublicKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "public key is not trusted")
	}

	verified := ed25519.Verify(pubKey, []byte(in.Message), in.Signature)
	if !verified {
		return nil, status.Errorf(codes.Unauthenticated, "signature is not valid")
	}

	return &pb.EchoReply{Message: "Echo " + in.Message}, nil
}

func NewServer() *Server {
	return &Server{}
}

func NewServerWithTrustedKeys(trustedKeys keystore.KeyStore) *Server {
	return &Server{trustedKeys: trustedKeys}
}

func (s *Server) Serve() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterEchoServiceServer(s.grpcServer, s)

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}
