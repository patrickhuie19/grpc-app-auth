package intgtest

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	client "grpc-app-auth/client"
	server "grpc-app-auth/server"
	"testing"
)

type TrustedKeyStore struct {
	keys map[string]ed25519.PublicKey
}

func (tks *TrustedKeyStore) GetPublicKey(keyID string) ([]byte, error) {
	key, ok := tks.keys[keyID]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return key, nil

}

func (tks *TrustedKeyStore) StorePublicKey(keyID string, publicKey []byte) error {
	tks.keys[keyID] = publicKey
	return nil
}

func Test(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	tks := &TrustedKeyStore{keys: map[string]ed25519.PublicKey{}}
	tks.StorePublicKey(base64.StdEncoding.EncodeToString(publicKey), publicKey)
	t.Logf("server public key string: %s", base64.StdEncoding.EncodeToString(publicKey))
	
	s := server.NewServerWithTrustedKeys(tks)
	c := client.NewClientWithKeys(publicKey, privateKey)

	go s.Serve()
	t.Cleanup(s.Stop)
	c.Echo("Hello World")
}
