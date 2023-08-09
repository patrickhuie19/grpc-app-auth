package keyutils

import (
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"path/filepath"
)

const (
	fixturesDir = "../fixtures"
)

func SaveKeysToFiles(pairs ...*FileKeyPair) error {
	// Ensure the directory exists
	if _, err := os.Stat(fixturesDir); os.IsNotExist(err) {
		err = os.Mkdir(fixturesDir, 0755)
		if err != nil {
			return err
		}
	}

	for _, pair := range pairs {
		// Define the path to the key file
		keyFile := filepath.Join(fixturesDir, pair.FileName)

		// Check if file already exists
		if _, err := os.Stat(keyFile); !os.IsNotExist(err) {
			continue
		}

		// Encode key in base64
		keyBase64 := base64.StdEncoding.EncodeToString(pair.Key)

		// Save key to file
		file, err := os.Create(keyFile)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(keyBase64); err != nil {
			return err
		}
	}

	return nil
}

func ReadKeysFromFiles(pubKeyFileName string, pubKey *ed25519.PublicKey, privKeyFileName string, privKey *ed25519.PrivateKey) error {
	// Read public key from file
	pubKeyFile := filepath.Join(fixturesDir, pubKeyFileName)
	pubKeyBase64, err := os.ReadFile(pubKeyFile)
	if err != nil {
		return err
	}

	// Read private key from file
	privKeyFile := filepath.Join(fixturesDir, privKeyFileName)
	privKeyBase64, err := os.ReadFile(privKeyFile)
	if err != nil {
		return err
	}

	// Decode public key from base64
	pubKeyDecoded, err := base64.StdEncoding.DecodeString(string(pubKeyBase64))
	if err != nil {
		return err
	}

	// Decode private key from base64
	privKeyDecoded, err := base64.StdEncoding.DecodeString(string(privKeyBase64))
	if err != nil {
		return err
	}

	ed25519pubKey := (ed25519.PublicKey)(pubKeyDecoded)
	ed25519privKey := (ed25519.PrivateKey)(privKeyDecoded)

	*pubKey = ed25519pubKey
	*privKey = ed25519privKey

	return nil
}

type FileKeyPair struct {
	FileName string
	Key      []byte
}
