#!/bin/bash

# Signature Service Demo Script
# This script demonstrates the API functionality

set -e

BASE_URL="http://localhost:8080/api/v0"

echo "=========================================="
echo "Signature Service API Demo"
echo "=========================================="
echo ""

# Check if server is running
echo "1. Checking server health..."
curl -s "$BASE_URL/health" | jq '.'
echo ""

# Create RSA device
echo "2. Creating RSA signature device..."
curl -s -X POST "$BASE_URL/devices" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "rsa-device-001",
    "algorithm": "RSA",
    "label": "RSA Demo Device"
  }' | jq '.'
echo ""

# Create ECC device
echo "3. Creating ECC signature device..."
curl -s -X POST "$BASE_URL/devices" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "ecc-device-001",
    "algorithm": "ECC",
    "label": "ECC Demo Device"
  }' | jq '.'
echo ""

# List all devices
echo "4. Listing all devices..."
curl -s "$BASE_URL/devices" | jq '.'
echo ""

# Get specific device
echo "5. Getting RSA device details..."
curl -s "$BASE_URL/devices/rsa-device-001" | jq '.'
echo ""

# Sign first transaction
echo "6. Signing first transaction (counter = 0)..."
curl -s -X POST "$BASE_URL/devices/rsa-device-001/sign" \
  -H "Content-Type: application/json" \
  -d '{
    "data": "transaction-001-purchase-100-USD"
  }' | jq '.'
echo ""

# Sign second transaction
echo "7. Signing second transaction (counter = 1)..."
curl -s -X POST "$BASE_URL/devices/rsa-device-001/sign" \
  -H "Content-Type: application/json" \
  -d '{
    "data": "transaction-002-refund-25-USD"
  }' | jq '.'
echo ""

# Sign third transaction
echo "8. Signing third transaction (counter = 2)..."
curl -s -X POST "$BASE_URL/devices/rsa-device-001/sign" \
  -H "Content-Type: application/json" \
  -d '{
    "data": "transaction-003-purchase-50-USD"
  }' | jq '.'
echo ""

# Verify counter incremented
echo "9. Verifying signature counter incremented..."
curl -s "$BASE_URL/devices/rsa-device-001" | jq '.data | {id, signature_counter, last_signature}'
echo ""

# Sign with ECC device
echo "10. Signing with ECC device..."
curl -s -X POST "$BASE_URL/devices/ecc-device-001/sign" \
  -H "Content-Type: application/json" \
  -d '{
    "data": "ecc-transaction-001"
  }' | jq '.'
echo ""

# Error case: Invalid algorithm
echo "11. Testing error case - invalid algorithm..."
curl -s -X POST "$BASE_URL/devices" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "invalid-device",
    "algorithm": "INVALID",
    "label": "Invalid Device"
  }' | jq '.'
echo ""

# Error case: Duplicate device
echo "12. Testing error case - duplicate device ID..."
curl -s -X POST "$BASE_URL/devices" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "rsa-device-001",
    "algorithm": "RSA",
    "label": "Duplicate Device"
  }' | jq '.'
echo ""

# Error case: Device not found
echo "13. Testing error case - device not found..."
curl -s "$BASE_URL/devices/nonexistent-device" | jq '.'
echo ""

echo "=========================================="
echo "Demo completed successfully!"
echo "=========================================="
