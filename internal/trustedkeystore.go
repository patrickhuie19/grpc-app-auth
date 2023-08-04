package internal

import (
	"crypto/ed25519"
	"fmt"
)

type TrustedKeyStore struct {
	Keys map[string]ed25519.PublicKey
}

func (tks *TrustedKeyStore) GetPublicKey(keyID string) ([]byte, error) {
	key, ok := tks.Keys[keyID]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return key, nil

}

func (tks *TrustedKeyStore) StorePublicKey(keyID string, publicKey []byte) error {
	tks.Keys[keyID] = publicKey
	return nil
}
