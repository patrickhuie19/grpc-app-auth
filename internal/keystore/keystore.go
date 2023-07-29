package keystore


type KeyStore interface {
	// GetPublicKey returns the public key for the given key ID.
	GetPublicKey(keyID string) ([]byte, error)

	StorePublicKey(keyID string, publicKey []byte) error
}
