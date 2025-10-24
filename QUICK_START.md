# Quick Start Guide

## Get Started in 3 Steps

### Step 1: Build
```bash
go build -o signature-service .
```

### Step 2: Run
```bash
./signature-service
```

You should see:
```
2024/10/18 21:00:00 Starting signature service on :8080
```

### Step 3: Test
Open a new terminal and run:
```bash
./demo.sh
```

---

## Manual Testing

### Create a Device
```bash
curl -X POST http://localhost:8080/api/v0/devices \
  -H "Content-Type: application/json" \
  -d '{
    "id": "my-device-001",
    "algorithm": "RSA",
    "label": "My First Device"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "id": "my-device-001",
    "label": "My First Device",
    "algorithm": "RSA",
    "signature_counter": 0,
    "last_signature": ""
  }
}
```

### Sign Your First Transaction
```bash
curl -X POST http://localhost:8080/api/v0/devices/my-device-001/sign \
  -H "Content-Type: application/json" \
  -d '{
    "data": "purchase-100-USD"
  }' | jq
```

**Expected Response:**
```json
{
  "data": {
    "signature": "VGhpcyBpcyBhIGJhc2U2NCBlbmNvZGVkIHNpZ25hdHVyZQ==",
    "signed_data": "0_purchase-100-USD_bXktZGV2aWNlLTAwMQ=="
  }
}
```

### Get Device Details
```bash
curl http://localhost:8080/api/v0/devices/my-device-001 | jq
```

**Notice:** `signature_counter` is now `1` and `last_signature` is populated!

### List All Devices
```bash
curl http://localhost:8080/api/v0/devices | jq
```

---

## Run Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./domain
go test -v ./api
```

---

## Key Concepts

### Signature Chain
Each signature includes:
1. **Counter**: Current signature count
2. **Data**: Your transaction data
3. **Last Signature**: Previous signature (or device ID for first)

**Format:** `<counter>_<data>_<last_signature_base64>`

### Example Chain
```
Signature 1: 0_purchase-100_base64(device-id)
Signature 2: 1_refund-25_base64(signature-1)
Signature 3: 2_purchase-50_base64(signature-2)
```

### Algorithms
- **RSA**: Traditional, widely supported
- **ECC**: Faster, smaller signatures

---

## Monitoring

### Health Check
```bash
curl http://localhost:8080/api/v0/health
```

Should return:
```json
{
  "data": {
    "status": "pass",
    "version": "v0"
  }
}
```

---

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill it
kill -9 <PID>

# Or change port in main.go
const ListenAddress = ":8081"
```

### Tests Failing
```bash
# Clean and rebuild
go clean
go build -o signature-service .
go test -v ./...
```

### Demo Script Not Executable
```bash
chmod +x demo.sh
```

---

## Next Steps

1. **Read API Documentation**: See `API_DOCUMENTATION.md`
2. **Understand Architecture**: See `IMPLEMENTATION.md`
3. **Explore Tests**: Check `*_test.go` files
4. **Experiment**: Try different algorithms, concurrent requests

---

## Pro Tips

### Pretty Print JSON
```bash
# Install jq if not available
sudo apt-get install jq  # Ubuntu/Debian
brew install jq          # macOS

# Use it with curl
curl http://localhost:8080/api/v0/devices | jq '.'
```

### Test Concurrent Access
```bash
# Run multiple signs in parallel
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v0/devices/my-device-001/sign \
    -H "Content-Type: application/json" \
    -d "{\"data\": \"transaction-$i\"}" &
done
wait

# Check final counter (should be 10)
curl http://localhost:8080/api/v0/devices/my-device-001 | jq '.data.signature_counter'
```

### Compare Algorithms
```bash
# Create RSA device
curl -X POST http://localhost:8080/api/v0/devices \
  -d '{"id": "rsa-test", "algorithm": "RSA"}' | jq

# Create ECC device
curl -X POST http://localhost:8080/api/v0/devices \
  -d '{"id": "ecc-test", "algorithm": "ECC"}' | jq

# Sign with both and compare signature lengths
curl -X POST http://localhost:8080/api/v0/devices/rsa-test/sign \
  -d '{"data": "test"}' | jq '.data.signature' | wc -c

curl -X POST http://localhost:8080/api/v0/devices/ecc-test/sign \
  -d '{"data": "test"}' | jq '.data.signature' | wc -c
```

---

## Learning Resources

- **Go Concurrency**: Understanding `sync.RWMutex` and `sync.Map`
- **Cryptography**: RSA-PSS vs ECDSA
- **REST API Design**: Resource-oriented architecture
- **Clean Architecture**: Separation of concerns

Happy coding! 🚀
