package intgtest

import (
	"crypto/ed25519"
	"encoding/base64"
	client "grpc-app-auth/client"
	"grpc-app-auth/internal"
	server "grpc-app-auth/server"
	"testing"
)

func Test(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	tks := &internal.TrustedKeyStore{Keys: map[string]ed25519.PublicKey{}}
	tks.StorePublicKey(base64.StdEncoding.EncodeToString(publicKey), publicKey)

	s := server.NewServerWithTrustedKeys(tks)
	c := client.NewClientWithKeys(publicKey, privateKey)

	go s.Serve()
	t.Cleanup(s.Stop)
	c.Echo("Hello World")
}
