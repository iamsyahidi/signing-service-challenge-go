package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

var (
	// ErrDeviceNotFound is returned when a device doesn't exist
	ErrDeviceNotFound = errors.New("device not found")
	// ErrDeviceAlreadyExists is returned when trying to create a device with an existing ID
	ErrDeviceAlreadyExists = errors.New("device already exists")
)

// InMemoryDeviceRepository implements domain.DeviceRepository with thread-safe in-memory storage
type InMemoryDeviceRepository struct {
	mu      sync.RWMutex
	devices map[string]*domain.SignatureDevice
}

// NewInMemoryDeviceRepository creates a new in-memory device repository
func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// Create adds a new device to the repository
func (r *InMemoryDeviceRepository) Create(device *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return ErrDeviceAlreadyExists
	}

	r.devices[device.ID] = device
	return nil
}

// GetByID retrieves a device by its ID
func (r *InMemoryDeviceRepository) GetByID(id string) (*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, ErrDeviceNotFound
	}

	return device, nil
}

// Update updates an existing device in the repository
func (r *InMemoryDeviceRepository) Update(device *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; !exists {
		return ErrDeviceNotFound
	}

	r.devices[device.ID] = device
	return nil
}

// List returns all devices in the repository
func (r *InMemoryDeviceRepository) List() ([]*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*domain.SignatureDevice, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}

	return devices, nil
}
