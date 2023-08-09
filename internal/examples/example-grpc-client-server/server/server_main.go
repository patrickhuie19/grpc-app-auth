package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"grpc-app-auth/internal"
	"grpc-app-auth/internal/keyutils"
	"grpc-app-auth/server"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var pubKey ed25519.PublicKey
	var privKey ed25519.PrivateKey
	err := keyutils.ReadKeysFromFiles("public.key", &pubKey, "private.key", &privKey)
	if err != nil {
		publicKey, privateKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatalf("Error generating keys: %v", err)
		}

		publicKeyPair := &keyutils.FileKeyPair{FileName: "public.key", Key: publicKey}
		privateKeyPair := &keyutils.FileKeyPair{FileName: "private.key", Key: privateKey}

		if err := keyutils.SaveKeysToFiles(publicKeyPair, privateKeyPair); err != nil {
			log.Fatalf("Error writing keys: %v", err)
		}
		pubKey = publicKey
	}

	tks := &internal.TrustedKeyStore{Keys: map[string]ed25519.PublicKey{}}
	tks.StorePublicKey(base64.StdEncoding.EncodeToString(pubKey), pubKey)

	enableTelemetry := os.Getenv("ENABLE_TELEMETRY")
	telemetryTarget := os.Getenv("TELEMETRY_TARGET")

	opts := make([]server.ServerOption, 0, 1)
	if enableTelemetry == "true" {
		opts = append(opts, server.WithOpenTelemetry(telemetryTarget))
	}

	s, err := server.NewServerWithTrustedKeysAndFuncOpts(tks, opts...)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Serve()
	}()
	time.AfterFunc(300*time.Second, s.Stop)
	wg.Wait()
}
