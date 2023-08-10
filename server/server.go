package server

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-app-auth/services"
	"grpc-app-auth/utils"

	"grpc-app-auth/internal/keystore"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedEchoServer
	pb.UnimplementedAddServer
	trustedKeys    keystore.KeyStore
	grpcServer     *grpc.Server
	tracerProvider *sdktrace.TracerProvider
}

type ServerOption func(*serverOptions) error

type serverOptions struct {
	enableTracing bool
	tracingTarget string
}

func WithOpenTelemetry(target string) ServerOption {
	return func(o *serverOptions) error {
		o.enableTracing = true
		o.tracingTarget = target
		return nil
	}
}

func NewServerWithTrustedKeys(trustedKeys keystore.KeyStore) *Server {
	return &Server{trustedKeys: trustedKeys}
}

func NewServerWithTrustedKeysAndFuncOpts(trustedKeys keystore.KeyStore, opts ...ServerOption) (*Server, error) {
	server := NewServerWithTrustedKeys(trustedKeys)

	// apply defaults
	o := &serverOptions{}

	// apply user options
	for _, opt := range opts {
		err := opt(o)
		if err != nil {
			return nil, fmt.Errorf("error applying option: %w", err)
		}
	}

	if o.enableTracing {
		if err := server.SetupOpenTelemetry(o.tracingTarget, "grpc-app-auth-server"); err != nil {
			return nil, err
		}
	}

	return server, nil
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

func (s *Server) Serve() {
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

func (s *Server) SetupOpenTelemetry(target string, serviceName string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)

	// Enough to shutdown the underlying connection since DialContext is used in blocking mode
	defer cancel()
	conn, err := grpc.DialContext(ctx, target,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Print(fmt.Errorf("failed to create gRPC connection to collector: %w", err))
		return err
	}

	// Set up a trace exporter
	// Shuttting down the traceExporter will not shutdown the underlying connection.
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		fmt.Print(fmt.Errorf("failed to create trace exporter: %w", err))
		return err
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion("0.0.1"),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	s.tracerProvider = tracerProvider
	return nil
}
