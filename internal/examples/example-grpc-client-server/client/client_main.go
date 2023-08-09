package main

import (
	"crypto/ed25519"
	"grpc-app-auth/client"
	"grpc-app-auth/internal/keyutils"
)

func main() {
	var pubKey ed25519.PublicKey
	var privKey ed25519.PrivateKey
	err := keyutils.ReadKeysFromFiles("public.key", &pubKey, "private.key", &privKey)
	if err != nil {
		panic(err)
	}

	c := client.NewClientWithKeys(pubKey, privKey)

	c.Echo("Hello World")

	c.Add(1, 2)
}
