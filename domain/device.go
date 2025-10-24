package domain

import (
	"encoding/base64"
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// SignatureAlgorithm represents the supported signing algorithms
type SignatureAlgorithm string

const (
	AlgorithmRSA SignatureAlgorithm = "RSA"
	AlgorithmECC SignatureAlgorithm = "ECC"
)

// SignatureDevice represents a device that can create cryptographic signatures
type SignatureDevice struct {
	ID               string             `json:"id"`
	Label            string             `json:"label"`
	Algorithm        SignatureAlgorithm `json:"algorithm"`
	SignatureCounter int                `json:"signature_counter"`
	LastSignature    string             `json:"last_signature"` // base64 encoded
	Signer           crypto.Signer      `json:"-"`              // Not serialized
}

// SignatureResponse represents the result of a signing operation
type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

// CreateSignatureDeviceRequest represents the request to create a new device
type CreateSignatureDeviceRequest struct {
	ID        string             `json:"id"`
	Algorithm SignatureAlgorithm `json:"algorithm"`
	Label     string             `json:"label,omitempty"`
}

// SignTransactionRequest represents the request to sign data
type SignTransactionRequest struct {
	Data string `json:"data"`
}

// IsValidAlgorithm checks if the algorithm is supported
func IsValidAlgorithm(algo SignatureAlgorithm) bool {
	return algo == AlgorithmRSA || algo == AlgorithmECC
}

// GetSecuredDataToSign constructs the data string to be signed according to spec
// Format: <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
func (d *SignatureDevice) GetSecuredDataToSign(dataToBeSigned string) string {
	lastSig := d.LastSignature
	if d.SignatureCounter == 0 {
		// Base case: use base64 encoded device ID
		lastSig = base64.StdEncoding.EncodeToString([]byte(d.ID))
	}
	return fmt.Sprintf("%d_%s_%s", d.SignatureCounter, dataToBeSigned, lastSig)
}

// Sign performs the signing operation and updates the device state
func (d *SignatureDevice) Sign(dataToBeSigned string) (*SignatureResponse, error) {
	securedData := d.GetSecuredDataToSign(dataToBeSigned)

	signature, err := d.Signer.Sign([]byte(securedData))
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	// Update device state
	d.LastSignature = base64.StdEncoding.EncodeToString(signature)
	d.SignatureCounter++

	return &SignatureResponse{
		Signature:  d.LastSignature,
		SignedData: securedData,
	}, nil
}
