package domain

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

// DeviceService handles business logic for signature devices
type DeviceService struct {
	repository DeviceRepository
	// deviceLocks provides per-device mutex for signing operations
	// This ensures signature_counter is strictly monotonic without gaps
	deviceLocks sync.Map
}

// NewDeviceService creates a new device service
func NewDeviceService(repository DeviceRepository) *DeviceService {
	return &DeviceService{
		repository: repository,
	}
}

// CreateDevice creates a new signature device with the specified algorithm
func (s *DeviceService) CreateDevice(req *CreateSignatureDeviceRequest) (*SignatureDevice, error) {
	// Validate algorithm
	if !IsValidAlgorithm(req.Algorithm) {
		return nil, fmt.Errorf("unsupported algorithm: %s", req.Algorithm)
	}

	// Validate ID
	if req.ID == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	// Generate key pair and create signer based on algorithm
	var signer crypto.Signer
	var err error

	switch req.Algorithm {
	case AlgorithmRSA:
		generator := new(crypto.RSAGenerator)
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
		}
		signer = crypto.NewRSASigner(keyPair.Private)

	case AlgorithmECC:
		generator := new(crypto.ECCGenerator)
		keyPair, err := generator.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to generate ECC key pair: %w", err)
		}
		signer = crypto.NewECDSASigner(keyPair.Private)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", req.Algorithm)
	}

	// Create device
	device := &SignatureDevice{
		ID:               req.ID,
		Label:            req.Label,
		Algorithm:        req.Algorithm,
		SignatureCounter: 0,
		LastSignature:    "",
		Signer:           signer,
	}

	// Persist device
	if err = s.repository.Create(device); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return device, nil
}

// GetDevice retrieves a device by ID
func (s *DeviceService) GetDevice(id string) (*SignatureDevice, error) {
	return s.repository.GetByID(id)
}

// ListDevices returns all devices
func (s *DeviceService) ListDevices() ([]*SignatureDevice, error) {
	return s.repository.List()
}

// SignTransaction signs data with the specified device
// This method is thread-safe and ensures the signature counter is monotonic
func (s *DeviceService) SignTransaction(deviceID string, dataToBeSigned string) (*SignatureResponse, error) {
	// Validate input
	if dataToBeSigned == "" {
		return nil, fmt.Errorf("data to be signed cannot be empty")
	}

	// Get or create device-specific lock
	lockInterface, _ := s.deviceLocks.LoadOrStore(deviceID, &sync.Mutex{})
	deviceLock := lockInterface.(*sync.Mutex)

	// Lock this specific device for the signing operation
	// This ensures signature_counter increments are atomic and sequential
	deviceLock.Lock()
	defer deviceLock.Unlock()

	// Retrieve device
	device, err := s.repository.GetByID(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device: %w", err)
	}

	// Perform signing (this updates device state)
	response, err := device.Sign(dataToBeSigned)
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	// Persist updated device state
	if err := s.repository.Update(device); err != nil {
		return nil, fmt.Errorf("failed to update device state: %w", err)
	}

	return response, nil
}
