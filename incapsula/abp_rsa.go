package incapsula

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func encryptRsa(rsaPem []byte, plaintext []byte, encLabel string) (string, error) {
	// Just taking the first entity here
	block, _ := pem.Decode(rsaPem)
	if block == nil {
		return "", fmt.Errorf("error decoding RSA PEM no entities found")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	var rsaKey = key.(*rsa.PublicKey)

	if err != nil {
		return "", fmt.Errorf("error parsing RSA public key: %w", err)
	}
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, []byte(encLabel))
	if err != nil {
		return "", fmt.Errorf("error encrypting value with RSA: %w", err)
	}
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	return encryptedBase64, nil
}
