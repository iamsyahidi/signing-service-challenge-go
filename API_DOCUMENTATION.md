# Signature Service API Documentation

## Overview

The Signature Service provides a RESTful API for managing cryptographic signature devices and signing transactions.

## Base URL

```
http://localhost:8080/api/v0
```

## Endpoints

### 1. Health Check

**GET** `/health`

Check the health status of the service.

**Response:**
```json
{
  "data": {
    "status": "pass",
    "version": "v0"
  }
}
```

---

### 2. Create Signature Device

**POST** `/devices`

Create a new signature device with a specific cryptographic algorithm.

**Request Body:**
```json
{
  "id": "device-123",
  "algorithm": "RSA",
  "label": "My Signature Device"
}
```

**Parameters:**
- `id` (string, required): Unique identifier for the device
- `algorithm` (string, required): Signing algorithm - either "RSA" or "ECC"
- `label` (string, optional): Human-readable label for the device

**Response (201 Created):**
```json
{
  "data": {
    "id": "device-123",
    "label": "My Signature Device",
    "algorithm": "RSA",
    "signature_counter": 0,
    "last_signature": ""
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid algorithm or missing required fields
- `409 Conflict`: Device with the same ID already exists

---

### 3. Get Signature Device

**GET** `/devices/{id}`

Retrieve details of a specific signature device.

**Response (200 OK):**
```json
{
  "data": {
    "id": "device-123",
    "label": "My Signature Device",
    "algorithm": "RSA",
    "signature_counter": 5,
    "last_signature": "base64-encoded-signature"
  }
}
```

**Error Responses:**
- `404 Not Found`: Device does not exist

---

### 4. List All Devices

**GET** `/devices`

Retrieve a list of all signature devices.

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "device-123",
      "label": "My Signature Device",
      "algorithm": "RSA",
      "signature_counter": 5,
      "last_signature": "base64-encoded-signature"
    },
    {
      "id": "device-456",
      "label": "Another Device",
      "algorithm": "ECC",
      "signature_counter": 10,
      "last_signature": "base64-encoded-signature"
    }
  ]
}
```

---

### 5. Sign Transaction

**POST** `/devices/{id}/sign`

Sign transaction data using the specified device.

**Request Body:**
```json
{
  "data": "transaction-data-to-sign"
}
```

**Parameters:**
- `data` (string, required): The data to be signed

**Response (200 OK):**
```json
{
  "data": {
    "signature": "base64-encoded-signature",
    "signed_data": "5_transaction-data-to-sign_previous-signature-base64"
  }
}
```

**Signed Data Format:**
```
<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
```

For the first signature (counter = 0), the last signature is replaced with the base64-encoded device ID.

**Error Responses:**
- `404 Not Found`: Device does not exist
- `400 Bad Request`: Missing or invalid data
- `500 Internal Server Error`: Signing operation failed

---

## Signature Chain Verification

The signature service implements a chain of trust:

1. **First Signature**: Uses base64(device_id) as the initial "last signature"
2. **Subsequent Signatures**: Each signature includes the previous signature in the data being signed
3. **Monotonic Counter**: The signature counter strictly increments without gaps
4. **Thread-Safe**: Concurrent signing operations are properly synchronized

### Example Signature Chain

```
Transaction 1:
  Input: "purchase-100"
  Secured Data: "0_purchase-100_ZGV2aWNlLTEyMw=="
  Signature: "sig1..."
  Counter: 0 → 1

Transaction 2:
  Input: "refund-50"
  Secured Data: "1_refund-50_sig1..."
  Signature: "sig2..."
  Counter: 1 → 2
```

---

## Error Response Format

All error responses follow this structure:

```json
{
  "errors": [
    "Error message 1",
    "Error message 2"
  ]
}
```

---

## Supported Algorithms

### RSA
- Key size: 512 bits (for demo purposes)
- Signature scheme: RSA-PSS with SHA256
- Suitable for: High security requirements

### ECC (Elliptic Curve Cryptography)
- Curve: P-384
- Signature scheme: ECDSA with SHA256
- Suitable for: Performance-critical applications

---

## Thread Safety

The service is designed for concurrent access:

- **Repository Level**: Read-write mutex protects the in-memory storage
- **Device Level**: Per-device mutex ensures atomic signing operations
- **Counter Guarantee**: Signature counters are strictly monotonic without gaps

---

## Testing

Run the test suite:

```bash
go test -v ./...
```

Tests cover:
- Device creation and retrieval
- Signature generation and verification
- Monotonic counter behavior
- Concurrent access scenarios
- Error handling
