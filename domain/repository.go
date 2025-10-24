package domain

// DeviceRepository defines the interface for device persistence operations
type DeviceRepository interface {
	Create(device *SignatureDevice) error
	GetByID(id string) (*SignatureDevice, error)
	Update(device *SignatureDevice) error
	List() ([]*SignatureDevice, error)
}
