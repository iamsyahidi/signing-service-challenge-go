package domain

import (
	"encoding/base64"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestGetSecuredDataToSign_FirstSignature(t *testing.T) {
	// Generate a test key pair
	generator := &crypto.RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	device := &SignatureDevice{
		ID:               "test-device-123",
		Label:            "Test Device",
		Algorithm:        AlgorithmRSA,
		SignatureCounter: 0,
		LastSignature:    "",
		Signer:           crypto.NewRSASigner(keyPair.Private),
	}

	data := "transaction-data"
	securedData := device.GetSecuredDataToSign(data)

	expectedLastSig := base64.StdEncoding.EncodeToString([]byte(device.ID))
	expected := "0_transaction-data_" + expectedLastSig

	if securedData != expected {
		t.Errorf("Expected %s, got %s", expected, securedData)
	}
}

func TestGetSecuredDataToSign_SubsequentSignature(t *testing.T) {
	generator := &crypto.RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	device := &SignatureDevice{
		ID:               "test-device-123",
		Label:            "Test Device",
		Algorithm:        AlgorithmRSA,
		SignatureCounter: 5,
		LastSignature:    "previous-signature-base64",
		Signer:           crypto.NewRSASigner(keyPair.Private),
	}

	data := "transaction-data"
	securedData := device.GetSecuredDataToSign(data)

	expected := "5_transaction-data_previous-signature-base64"

	if securedData != expected {
		t.Errorf("Expected %s, got %s", expected, securedData)
	}
}

func TestSign_UpdatesDeviceState(t *testing.T) {
	generator := &crypto.RSAGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	device := &SignatureDevice{
		ID:               "test-device-123",
		Label:            "Test Device",
		Algorithm:        AlgorithmRSA,
		SignatureCounter: 0,
		LastSignature:    "",
		Signer:           crypto.NewRSASigner(keyPair.Private),
	}

	initialCounter := device.SignatureCounter

	response, err := device.Sign("test-data")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	// Verify counter incremented
	if device.SignatureCounter != initialCounter+1 {
		t.Errorf("Expected counter to be %d, got %d", initialCounter+1, device.SignatureCounter)
	}

	// Verify last signature was updated
	if device.LastSignature == "" {
		t.Error("LastSignature should not be empty after signing")
	}

	// Verify response contains signature
	if response.Signature == "" {
		t.Error("Response signature should not be empty")
	}

	// Verify response contains signed data
	if response.SignedData == "" {
		t.Error("Response signed_data should not be empty")
	}

	// Verify signature matches last signature
	if response.Signature != device.LastSignature {
		t.Error("Response signature should match device LastSignature")
	}
}

func TestSign_MonotonicCounter(t *testing.T) {
	generator := &crypto.ECCGenerator{}
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	device := &SignatureDevice{
		ID:               "test-device-456",
		Label:            "Test Device",
		Algorithm:        AlgorithmECC,
		SignatureCounter: 0,
		LastSignature:    "",
		Signer:           crypto.NewECDSASigner(keyPair.Private),
	}

	// Perform multiple signatures
	for i := 0; i < 10; i++ {
		expectedCounter := i
		if device.SignatureCounter != expectedCounter {
			t.Errorf("Expected counter %d, got %d", expectedCounter, device.SignatureCounter)
		}

		_, err := device.Sign("data-" + string(rune(i)))
		if err != nil {
			t.Fatalf("Sign %d failed: %v", i, err)
		}
	}

	// Final counter should be 10
	if device.SignatureCounter != 10 {
		t.Errorf("Expected final counter to be 10, got %d", device.SignatureCounter)
	}
}

func TestIsValidAlgorithm(t *testing.T) {
	tests := []struct {
		algo  SignatureAlgorithm
		valid bool
	}{
		{AlgorithmRSA, true},
		{AlgorithmECC, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		result := IsValidAlgorithm(tt.algo)
		if result != tt.valid {
			t.Errorf("IsValidAlgorithm(%s) = %v, want %v", tt.algo, result, tt.valid)
		}
	}
}
