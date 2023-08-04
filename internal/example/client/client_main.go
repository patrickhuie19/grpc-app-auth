package main

import (
	"grpc-app-auth/client"
	"grpc-app-auth/internal/keyutils"
)

func main() {
	pubKey, privKey := keyutils.ReadKeysFromFiles("public.key", "private.key")

	c := client.NewClientWithKeys(pubKey, privKey)

	c.Echo("Hello World")
}
