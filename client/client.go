package client

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"time"

	pb "grpc-app-auth/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	publicKey ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func NewClient() *Client {
	return &Client{}
}

func NewClientWithKeys(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) *Client {
	return &Client{publicKey: publicKey, privateKey: privateKey}
}

func (c *Client) Echo(message string) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewEchoServiceClient(conn)


	signature := ed25519.Sign(c.privateKey, []byte(message))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pubKeyStr := base64.StdEncoding.EncodeToString(c.publicKey)
	log.Printf("client public key string: %s", base64.StdEncoding.EncodeToString(c.publicKey))
	r, err := grpcClient.Echo(ctx, &pb.EchoRequest{Message: message, PublicKey: pubKeyStr, Signature: signature})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
