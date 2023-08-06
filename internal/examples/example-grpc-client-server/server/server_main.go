package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"grpc-app-auth/internal"
	"grpc-app-auth/internal/keyutils"
	"grpc-app-auth/server"
	"log"
	"sync"
	"time"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	publicKeyPair := &keyutils.FileKeyPair{FileName: "public.key", Key: publicKey}
	privateKeyPair := &keyutils.FileKeyPair{FileName: "private.key", Key: privateKey}

	if err := keyutils.SaveKeysToFiles(publicKeyPair, privateKeyPair); err != nil {
		log.Fatalf("Error writing keys: %v", err)
	}

	tks := &internal.TrustedKeyStore{Keys: map[string]ed25519.PublicKey{}}
	tks.StorePublicKey(base64.StdEncoding.EncodeToString(publicKey), publicKey)

	s := server.NewServerWithTrustedKeys(tks)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Serve()
	}()
	time.AfterFunc(300*time.Second, s.Stop)
	wg.Wait()
}
