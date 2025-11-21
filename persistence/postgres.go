package persistence

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDeviceRepository implements domain.DeviceRepository with PostgreSQL storage
type PostgresDeviceRepository struct {
	db *sqlx.DB
}

// DeviceDB represents the database model for a signature device
type DeviceDB struct {
	ID               string    `db:"id"`
	Label           *string   `db:"label"`
	Algorithm       string    `db:"algorithm"`
	SignatureCounter int       `db:"signature_counter"`
	LastSignature   *string   `db:"last_signature"`
	PublicKey       string    `db:"public_key"`
	PrivateKey      string    `db:"private_key"`
	KeyType         string    `db:"key_type"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// NewPostgresDeviceRepository creates a new PostgreSQL device repository with a connection string
func NewPostgresDeviceRepository(connectionString string) (*PostgresDeviceRepository, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresDeviceRepository{
		db: db,
	}, nil
}

// NewPostgresDeviceRepositoryFromDB creates a new PostgreSQL device repository from an existing *sql.DB
func NewPostgresDeviceRepositoryFromDB(db *sql.DB) *PostgresDeviceRepository {
	return &PostgresDeviceRepository{
		db: sqlx.NewDb(db, "postgres"),
	}
}

// Create adds a new device to the repository
func (r *PostgresDeviceRepository) Create(device *domain.SignatureDevice) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if device already exists
	var exists bool
	err = tx.Get(&exists, `SELECT EXISTS(SELECT 1 FROM signature_devices WHERE id = $1)`, device.ID)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %w", err)
	}

	if exists {
		return ErrDeviceAlreadyExists
	}

	// Extract key information
	var publicKey, privateKey, keyType string

	// Get the private key using the GetPrivateKey method
	privKey := device.Signer.GetPrivateKey()

	switch key := privKey.(type) {
	case *rsa.PrivateKey:
		// Export private key as PEM
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
		privateKey = string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}))

		// Export public key as PEM
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal RSA public key: %w", err)
		}
		publicKey = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}))
		keyType = "RSA"

	case *ecdsa.PrivateKey:
		// Export private key as PEM
		privateKeyBytes, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return fmt.Errorf("failed to marshal ECDSA private key: %w", err)
		}
		privateKey = string(pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		}))

		// Export public key as PEM
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal ECDSA public key: %w", err)
		}
		publicKey = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}))
		keyType = "ECC"

	default:
		return fmt.Errorf("unsupported key type: %T", key)
	}

	// Insert the new device
	_, err = tx.Exec(`
		INSERT INTO signature_devices (
			id, label, algorithm, signature_counter, last_signature,
			public_key, private_key, key_type, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, 
		device.ID,
		device.Label,
		string(device.Algorithm),
		device.SignatureCounter,
		device.LastSignature,
		publicKey,
		privateKey,
		keyType,
	)

	if err != nil {
		return fmt.Errorf("failed to insert device: %w", err)
	}

	return tx.Commit()
}

// GetByID retrieves a device by its ID
func (r *PostgresDeviceRepository) GetByID(id string) (*domain.SignatureDevice, error) {
	var deviceDB DeviceDB
	err := r.db.Get(&deviceDB, `
		SELECT * FROM signature_devices WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return r.toDomainDevice(&deviceDB)
}

// Update updates an existing device in the repository
func (r *PostgresDeviceRepository) Update(device *domain.SignatureDevice) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if device exists
	var exists bool
	err = tx.Get(&exists, `SELECT EXISTS(SELECT 1 FROM signature_devices WHERE id = $1)`, device.ID)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %w", err)
	}

	if !exists {
		return ErrDeviceNotFound
	}

	// Update the device
	_, err = tx.Exec(`
		UPDATE signature_devices 
		SET 
			signature_counter = $1,
			last_signature = $2,
			updated_at = NOW()
		WHERE id = $3
	`, 
		device.SignatureCounter,
		device.LastSignature,
		device.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	return tx.Commit()
}

// List returns all devices in the repository
func (r *PostgresDeviceRepository) List() ([]*domain.SignatureDevice, error) {
	var devicesDB []DeviceDB
	err := r.db.Select(&devicesDB, `SELECT * FROM signature_devices`)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	devices := make([]*domain.SignatureDevice, 0, len(devicesDB))
	for _, deviceDB := range devicesDB {
		device, err := r.toDomainDevice(&deviceDB)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// Close closes the database connection
func (r *PostgresDeviceRepository) Close() error {
	return r.db.Close()
}

// toDomainDevice converts a DeviceDB to a domain.SignatureDevice
func (r *PostgresDeviceRepository) toDomainDevice(deviceDB *DeviceDB) (*domain.SignatureDevice, error) {
	var signer crypto.Signer

	// Decode private key
	block, _ := pem.Decode([]byte(deviceDB.PrivateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM block")
	}

	switch deviceDB.KeyType {
	case "RSA":
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}
		signer = crypto.NewRSASigner(privKey)
	case "ECC":
		privKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ECDSA private key: %w", err)
		}
		signer = crypto.NewECDSASigner(privKey)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", deviceDB.KeyType)
	}

	var lastSignature string
	if deviceDB.LastSignature != nil {
		lastSignature = *deviceDB.LastSignature
	}

	var label string
	if deviceDB.Label != nil {
		label = *deviceDB.Label
	}

	return &domain.SignatureDevice{
		ID:               deviceDB.ID,
		Label:            label,
		Algorithm:        domain.SignatureAlgorithm(deviceDB.Algorithm),
		SignatureCounter: deviceDB.SignatureCounter,
		LastSignature:    lastSignature,
		Signer:           signer,
	}, nil
}
