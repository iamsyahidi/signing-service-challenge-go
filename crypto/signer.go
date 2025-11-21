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
	GetPrivateKey() interface{}
}

// RSASigner implements the Signer interface for RSA signatures
type RSASigner struct {
	PrivateKey *rsa.PrivateKey
}

// NewRSASigner creates a new RSA signer with the given private key
func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		PrivateKey: privateKey,
	}
}

// Sign signs the data using RSA-PSS with SHA256
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	
	signature, err := rsa.SignPSS(
		rand.Reader,
		s.PrivateKey,
		crypto.SHA256,
		hash[:],
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("RSA signing failed: %w", err)
	}
	
	return signature, nil
}

// GetPrivateKey returns the RSA private key
func (s *RSASigner) GetPrivateKey() interface{} {
	return s.PrivateKey
}

// ECDSASigner implements the Signer interface for ECDSA signatures
type ECDSASigner struct {
	PrivateKey *ecdsa.PrivateKey
}

// NewECDSASigner creates a new ECDSA signer with the given private key
func NewECDSASigner(privateKey *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{
		PrivateKey: privateKey,
	}
}

// Sign signs the data using ECDSA with SHA256
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.Sum256(dataToBeSigned)
	r, s2, err := ecdsa.Sign(rand.Reader, s.PrivateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("ECDSA signing failed: %w", err)
	}

	signature := append(r.Bytes(), s2.Bytes()...)
	return signature, nil
}

// GetPrivateKey returns the ECDSA private key
func (s *ECDSASigner) GetPrivateKey() interface{} {
	return s.PrivateKey
}
