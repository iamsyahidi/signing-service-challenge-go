# Signature Service - Coding Challenge

## Instructions

This challenge is part of the software engineering interview process at fiskaly.

If you see this challenge, you've passed the first round of interviews and are now at the second and last stage.

We would like you to attempt the challenge below. You will then be able to discuss your solution in the skill-fit interview with two of our colleagues from the development department.

The quality of your code is more important to us than the quantity.

### Project Setup

For the challenge, we provide you with:

- Go project containing the setup
- Basic API structure and functionality
- Encoding / decoding of different key types (only needed to serialize keys to a persistent storage)
- Key generation algorithms (ECC, RSA)
- Library to generate UUIDs, included in `go.mod`

You can use these things as a foundation, but you're also free to modify them as you see fit.

### Prerequisites & Tooling

- Golang (v1.20+)

### The Challenge

The goal is to implement an API service that allows customers to create `signature devices` with which they can sign arbitrary transaction data.

#### Domain Description

The `signature service` can manage multiple `signature devices`. Such a device is identified by a unique identifier (e.g. UUID). For now you can pretend there is only one user / organization using the system (e.g. a dedicated node for them), therefore you do not need to think about user management at all.

When creating the `signature device`, the client of the API has to choose the signature algorithm that the device will be using to sign transaction data. During the creation process, a new key pair (`public key` & `private key`) has to be generated and assigned to the device.

The `signature device` should also have a `label` that can be used to display it in the UI and a `signature_counter` that tracks how many signatures have been created with this device. The `label` is provided by the user. The `signature_counter` shall only be modified internally.

##### Signature Creation

For the signature creation, the client will have to provide `data_to_be_signed` through the API. In order to increase the security of the system, we will extend this raw data with the current `signature_counter` and the `last_signature`.

The resulting string (`secured_data_to_be_signed`) should follow this format: `<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>`

In the base case there is no `last_signature` (= `signature_counter == 0`). Use the `base64`-encoded device ID (`last_signature = base64(device.id)`) instead of the `last_signature`.

This special string will be signed (`Signer.sign(secured_data_to_be_signed)`) and the resulting signature (`base64` encoded) will be returned to the client. The signature response could look like this:

```json
{ 
    "signature": <signature_base64_encoded>,
    "signed_data": "<signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>"
}
```

After the signature has been created, the signature counter's value has to be incremented (`signature_counter += 1`).

#### API

For now we need to provide two main operations to our customers:

- `CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
- `SignTransaction(deviceId: string, data: string): SignatureResponse`

Think of how to expose these operations through a RESTful HTTP-based API.

In addition, `list / retrieval operations` for the resources generated in the previous operations should be made available to the customers.

#### QA / Testing

As we are in the business of compliance technology, we need to make sure that our implementation is verifiably correct. Think of an automatable way to assure the correctness (in this challenge: adherence to the specifications) of the system.

#### Technical Constraints & Considerations

- The system will be used by many concurrent clients accessing the same resources.
- The `signature_counter` has to be strictly monotonically increasing and ideally without any gaps.
- The system currently only supports `RSA` and `ECDSA` as signature algorithms. Try to design the signing mechanism in a way that allows easy extension to other algorithms without changing the core domain logic.
- For now it is enough to store signature devices in memory. Efficiency is not a priority for this. In the future we might want to scale out. As you design your storage logic, keep in mind that we may later want to switch to a relational database.

### AI Tools

The use of AI tools to aid completing the challenge is permitted, but you will need to be able to reason about the design and implementation choices made when you reach the interview stage. Furthermore, if you used any AI tools, you need to clearly state which tools were used for different parts of the challenge. Ensure that you document this inside the `README` for your repository, so that it is visible to the reviewers.

### Credits

This challenge is heavily influenced by the regulations for `KassenSichV` (Germany) as well as the `RKSV` (Austria) and our solutions for them.

---

## Implementation Summary

### Completed Features

All requirements have been successfully implemented:

1. **Signature Device Management**
   - Create devices with RSA or ECC algorithms
   - Retrieve individual devices
   - List all devices
   - Automatic key pair generation

2. **Transaction Signing**
   - Sign arbitrary transaction data
   - Signature chaining (each signature includes previous)
   - Monotonic counter (strictly increasing, no gaps)
   - Base64 encoding of signatures

3. **Thread Safety**
   - Repository-level read-write locks
   - Per-device locks for signing operations
   - Tested with concurrent access scenarios

4. **RESTful API**
   - Resource-oriented endpoints
   - Proper HTTP status codes
   - JSON request/response format
   - Error handling

5. **Testing & QA**
   - Comprehensive unit tests
   - Integration tests
   - Concurrent access tests
   - All tests passing ✅

### Architecture

**Clean Architecture** with clear separation:
- **API Layer**: HTTP handlers, routing
- **Service Layer**: Business logic orchestration
- **Domain Layer**: Core entities, business rules
- **Persistence Layer**: Repository pattern (in-memory)

**Design Patterns Used:**
- Repository Pattern (easy database migration)
- Strategy Pattern (extensible signing algorithms)
- Dependency Injection
- Interface Segregation

### Quick Start

```bash
# Build the application
go build -o signature-service .

# Run the server
./signature-service

# In another terminal, run tests
go test -v ./...

# Run the demo
./demo.sh
```

### Documentation

- **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)**: Complete API reference with examples
- **[IMPLEMENTATION.md](IMPLEMENTATION.md)**: Architecture decisions, design patterns, and technical details

### API Endpoints

```
GET    /api/v0/health                → Health check
POST   /api/v0/devices               → Create signature device
GET    /api/v0/devices               → List all devices
GET    /api/v0/devices/{id}          → Get specific device
POST   /api/v0/devices/{id}/sign     → Sign transaction
```

### Example Usage

```bash
# Create a device
curl -X POST http://localhost:8080/api/v0/devices \
  -H "Content-Type: application/json" \
  -d '{
    "id": "device-001",
    "algorithm": "RSA",
    "label": "My Device"
  }'

# Sign a transaction
curl -X POST http://localhost:8080/api/v0/devices/device-001/sign \
  -H "Content-Type: application/json" \
  -d '{
    "data": "transaction-data"
  }'
```

### Testing

All tests pass successfully:

```bash
$ go test -v ./...
=== RUN   TestCreateDevice_Success
--- PASS: TestCreateDevice_Success (0.01s)
=== RUN   TestSignTransaction_ConcurrentAccess
--- PASS: TestSignTransaction_ConcurrentAccess (0.00s)
...
PASS
ok      .../api     0.027s
ok      .../domain  0.014s
```

**Test Coverage:**
- Device creation and validation
- Signature generation and chaining
- Monotonic counter behavior
- Concurrent access (critical!)
- Error handling

### Security Features

1. **Signature Chain**: Each signature cryptographically linked to previous
2. **Tamper Evidence**: Breaking chain invalidates subsequent signatures
3. **Monotonic Counter**: Prevents replay attacks
4. **Algorithm Support**: RSA-PSS and ECDSA with SHA256

### Technical Highlights

**Thread Safety Implementation:**
- Two-level locking strategy
- Repository-level: `sync.RWMutex` for concurrent reads
- Device-level: `sync.Map` with per-device mutexes
- Ensures counter is strictly monotonic without gaps

**Extensibility:**
- Adding new algorithms: Just implement `Signer` interface
- Switching to database: Implement `DeviceRepository` interface
- No changes to core domain logic required

**Compliance:**
- Aligns with KassenSichV and RKSV requirements
- Signature chain for audit trail
- Tamper-evident design

### Future Enhancements

Ready for production with minimal changes:
1. **Database Persistence**: PostgreSQL/MySQL (repository pattern ready)
2. **Key Management**: HSM/KMS integration
3. **Authentication**: JWT/OAuth2
4. **HTTPS**: TLS 1.3
5. **Additional Algorithms**: Ed25519, post-quantum signatures
6. **Audit Logging**: Comprehensive compliance logs
7. **Rate Limiting**: Per-device rate limits

### Contact

For questions about implementation decisions or architecture, please refer to:
- **IMPLEMENTATION.md** for detailed technical documentation
- **API_DOCUMENTATION.md** for API usage examples
- Test files for usage patterns and edge cases
