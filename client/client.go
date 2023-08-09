package client

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"time"

	pb "grpc-app-auth/services"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"grpc-app-auth/utils"
)

type Client struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func NewClient() *Client {
	return &Client{}
}

func NewClientWithKeys(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) *Client {
	return &Client{publicKey: publicKey, privateKey: privateKey}
}

func (c *Client) Echo(message string) {
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				otelgrpc.UnaryClientInterceptor(),
				loggingUnaryClientInterceptor,
			),
		),
	)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewEchoClient(conn)

	signature := ed25519.Sign(c.privateKey, []byte(message))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pubKeyStr := base64.StdEncoding.EncodeToString(c.publicKey)
	r, err := grpcClient.Echo(ctx, &pb.EchoRequest{Message: message, PublicKey: pubKeyStr, Signature: signature})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func (c *Client) Add(a float64, b float64) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewAddClient(conn)

	signature := ed25519.Sign(c.privateKey, []byte(utils.AddCanonicalization(a, b)))

	signatureStr := base64.StdEncoding.EncodeToString(signature)
	pubKeyStr := base64.StdEncoding.EncodeToString(c.publicKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "signature", signatureStr, "key", pubKeyStr)

	r, err := grpcClient.Add(ctx, &pb.AddRequest{A: a, B: b})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Result: %v", r.Result)
}

func loggingUnaryClientInterceptor(
	ctx context.Context, method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	log.Printf("[client] Request: %+v", req)
	err := invoker(ctx, method, req, resp, cc, opts...)
	log.Printf("[client] Response: %+v", resp)
	return err
}
