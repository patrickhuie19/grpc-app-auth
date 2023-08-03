package server

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"net"
	"time"

	pb "grpc-app-auth/services"
	"grpc-app-auth/utils"

	"grpc-app-auth/internal/keystore"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedEchoServer
	pb.UnimplementedAddServer
	trustedKeys keystore.KeyStore
	grpcServer  *grpc.Server
	tracerProvider *sdktrace.TracerProvider
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

func (s *Server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing authentication metadata")
	}

	pubKey, err := s.trustedKeys.GetPublicKey(md["key"][0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "public key is not trusted")
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(md["signature"][0])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "malformed signature")
	}

	verified := ed25519.Verify(pubKey, []byte(utils.AddCanonicalization(in.A, in.B)), signatureBytes)
	if !verified {
		return nil, status.Errorf(codes.Unauthenticated, "signature is not valid")
	}

	return &pb.AddReply{Result: in.A + in.B}, nil
}

func NewServer() *Server {
	return &Server{}
}

func NewServerWithTrustedKeys(trustedKeys keystore.KeyStore) *Server {
	return &Server{trustedKeys: trustedKeys}
}

func (s *Server) Serve() {
	s.setupOpenTelemetry()

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				loggingUnaryServerInterceptor,
			),
		),
	)
	pb.RegisterEchoServer(s.grpcServer, s)
	pb.RegisterAddServer(s.grpcServer, s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) Stop() {
	if s.tracerProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := s.tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shut down TracerProvider: %v", err)
		}
	}
	s.grpcServer.Stop()
}

func loggingUnaryServerInterceptor(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("[server] Metadata: %v", md)
	log.Printf("[server] Trace ID: %s", spanContext.TraceID().String())
	log.Printf("[server] Span ID: %s", spanContext.SpanID().String())
	log.Printf("[server] Request: %+v", req)
	resp, err := handler(ctx, req)
	log.Printf("[server] Response: %+v", resp)
	return resp, err
}

func (s *Server) setupOpenTelemetry() {
	// Configure OpenTelemetry SDK and exporters here
	// You may configure stdout exporter, Jaeger exporter, or others
	// Here's an example of setting up the stdout exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("grpc-app-auth-server"),
		semconv.ServiceVersion("0.0.1"),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	s.tracerProvider = tracerProvider
}
