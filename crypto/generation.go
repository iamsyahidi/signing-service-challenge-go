package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// GenerateRSAKeyPair generates a new RSA key pair with 2048 bits
func GenerateRSAKeyPair() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}
	return privateKey, nil
}

// GenerateECCKeyPair generates a new ECC key pair using P-256 curve
func GenerateECCKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECC key pair: %w", err)
	}
	return privateKey, nil
}

// ExportPrivateKeyToPEM exports a private key to PEM format
func ExportPrivateKeyToPEM(privateKey interface{}) (string, error) {
	var pemBlock *pem.Block

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		}
	case *ecdsa.PrivateKey:
		keyBytes, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return "", fmt.Errorf("failed to marshal ECDSA private key: %w", err)
		}
		pemBlock = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		}
	default:
		return "", errors.New("unsupported private key type")
	}

	return string(pem.EncodeToMemory(pemBlock)), nil
}

// ExportPublicKeyToPEM exports a public key to PEM format
func ExportPublicKeyToPEM(publicKey interface{}) (string, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	return string(pem.EncodeToMemory(pemBlock)), nil
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECC key pair: %w", err)
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}
