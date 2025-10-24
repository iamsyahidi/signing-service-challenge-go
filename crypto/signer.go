package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// RSASigner implements the Signer interface for RSA signatures
type RSASigner struct {
	privateKey *rsa.PrivateKey
}

// NewRSASigner creates a new RSA signer with the given private key
func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		privateKey: privateKey,
	}
}

// Sign signs the data using RSA-PSS with SHA256
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	
	signature, err := rsa.SignPSS(
		rand.Reader,
		s.privateKey,
		crypto.SHA256,
		hash[:],
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("RSA signing failed: %w", err)
	}
	
	return signature, nil
}

// ECDSASigner implements the Signer interface for ECDSA signatures
type ECDSASigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewECDSASigner creates a new ECDSA signer with the given private key
func NewECDSASigner(privateKey *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{
		privateKey: privateKey,
	}
}

// Sign signs the data using ECDSA with SHA256
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	
	signature, err := ecdsa.SignASN1(rand.Reader, s.privateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("ECDSA signing failed: %w", err)
	}
	
	return signature, nil
}
