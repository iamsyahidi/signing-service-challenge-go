# Project Structure

```
signing-service-challenge-go/
│
├── api/                          # API Layer (HTTP Handlers)
│   ├── device.go                 # Device endpoints (Create, Get, List, Sign)
│   ├── device_test.go            # Integration tests for API
│   ├── health.go                 # Health check endpoint
│   └── server.go                 # HTTP server setup and routing
│
├── crypto/                       # Cryptography Layer
│   ├── ecdsa.go                  # ECDSA key marshaling/unmarshaling
│   ├── generation.go             # Key pair generators (RSA, ECC)
│   ├── rsa.go                    # RSA key marshaling/unmarshaling
│   └── signer.go                 # Signer interface + implementations
│
├── domain/                       # Domain Layer (Business Logic)
│   ├── device.go                 # SignatureDevice entity and DTOs
│   ├── device_test.go            # Unit tests for domain logic
│   ├── repository.go             # Repository interface definition
│   └── service.go                # DeviceService (business logic)
│
├── persistence/                  # Persistence Layer
│   └── inmemory.go               # In-memory repository implementation
│
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
│
├── README.md                     # Main documentation
├── QUICK_START.md                # Getting started guide
├── API_DOCUMENTATION.md          # Complete API reference
├── IMPLEMENTATION.md             # Architecture deep dive
├── PROJECT_STRUCTURE.md          # This file
│
└── demo.sh                       # Interactive demo script
```

---

## Package Overview

### `api` Package
**Purpose**: HTTP layer, request/response handling

**Key Files**:
- `server.go`: Server initialization, routing, middleware
- `device.go`: Device-related endpoints
- `health.go`: Health check endpoint
- `device_test.go`: API integration tests

**Responsibilities**:
- Parse HTTP requests
- Validate input
- Call service layer
- Format responses
- Handle HTTP errors

---

### `crypto` Package
**Purpose**: Cryptographic operations

**Key Files**:
- `signer.go`: Signer interface + RSA/ECDSA implementations
- `generation.go`: Key pair generation
- `rsa.go`: RSA key serialization
- `ecdsa.go`: ECDSA key serialization

**Responsibilities**:
- Generate key pairs
- Sign data
- Serialize/deserialize keys

**Note**: Marshalers provided for future database persistence

---

### `domain` Package
**Purpose**: Core business logic and entities

**Key Files**:
- `device.go`: SignatureDevice entity, DTOs, business rules
- `service.go`: DeviceService orchestration
- `repository.go`: Repository interface
- `device_test.go`: Unit tests

**Responsibilities**:
- Define domain entities
- Implement business rules
- Orchestrate operations
- Validate business constraints

**Key Principle**: No dependencies on infrastructure (API, DB)

---

### `persistence` Package
**Purpose**: Data storage implementations

**Key Files**:
- `inmemory.go`: In-memory repository with thread safety

**Responsibilities**:
- Store and retrieve devices
- Ensure data consistency
- Thread-safe operations

**Future**: Add `postgres.go`, `mongodb.go`, etc.

---

## Data Flow

### Creating a Device
```
HTTP POST /api/v0/devices
    ↓
api.CreateDevice (parse request)
    ↓
domain.DeviceService.CreateDevice (business logic)
    ↓
crypto.Generator.Generate (create keys)
    ↓
crypto.NewSigner (create signer)
    ↓
persistence.Repository.Create (store device)
    ↓
api.WriteAPIResponse (return JSON)
```

### Signing a Transaction
```
HTTP POST /api/v0/devices/{id}/sign
    ↓
api.SignTransaction (parse request)
    ↓
domain.DeviceService.SignTransaction (orchestrate)
    ↓
[LOCK device-specific mutex]
    ↓
persistence.Repository.GetByID (fetch device)
    ↓
domain.SignatureDevice.Sign (business logic)
    ↓
crypto.Signer.Sign (cryptographic operation)
    ↓
persistence.Repository.Update (save state)
    ↓
[UNLOCK device-specific mutex]
    ↓
api.WriteAPIResponse (return signature)
```

---

## Dependency Graph

```
main.go
  ↓
api.Server
  ↓
domain.DeviceService
  ↓
domain.DeviceRepository (interface)
  ↑
persistence.InMemoryDeviceRepository (implementation)

domain.SignatureDevice
  ↓
crypto.Signer (interface)
  ↑
crypto.RSASigner / crypto.ECDSASigner (implementations)
```

**Key**: Dependencies point inward (Dependency Inversion Principle)

---

## Design Principles Applied

### 1. Clean Architecture
- **API Layer**: Depends on Domain
- **Domain Layer**: Independent (core business logic)
- **Persistence Layer**: Depends on Domain (implements interface)

### 2. SOLID Principles
- **S**ingle Responsibility: Each package has one reason to change
- **O**pen/Closed: Extensible via interfaces (Signer, Repository)
- **L**iskov Substitution: All implementations honor contracts
- **I**nterface Segregation: Small, focused interfaces
- **D**ependency Inversion: Depend on abstractions, not concretions

### 3. Repository Pattern
- Domain defines interface
- Persistence implements interface
- Easy to swap implementations

### 4. Strategy Pattern
- Signer interface
- Multiple implementations (RSA, ECDSA)
- Easy to add new algorithms

---

## File Responsibilities

### `main.go`
- Application bootstrap
- Dependency injection
- Server startup

### `api/server.go`
- HTTP server configuration
- Route registration
- Middleware setup
- Response helpers

### `api/device.go`
- Device CRUD endpoints
- Sign transaction endpoint
- Request parsing
- Response formatting

### `domain/device.go`
- SignatureDevice entity
- Request/Response DTOs
- Business rules (GetSecuredDataToSign, Sign)
- Validation logic

### `domain/service.go`
- Business logic orchestration
- Device creation workflow
- Transaction signing workflow
- Concurrency control (per-device locks)

### `domain/repository.go`
- Repository interface definition
- Storage contract

### `persistence/inmemory.go`
- In-memory storage implementation
- Thread-safe operations (RWMutex)
- CRUD operations

### `crypto/signer.go`
- Signer interface
- RSASigner implementation
- ECDSASigner implementation

### `crypto/generation.go`
- RSA key pair generation
- ECC key pair generation

---

## Test Files

### `api/device_test.go`
**Type**: Integration tests

**Coverage**:
- API endpoint functionality
- HTTP status codes
- Request/response formats
- Error handling
- Concurrent access

**Key Test**: `TestSignTransaction_ConcurrentAccess`

### `domain/device_test.go`
**Type**: Unit tests

**Coverage**:
- Secured data format
- Signature counter logic
- Device state updates
- Algorithm validation

---

## Documentation Files

### `README.md`
- Challenge description
- Implementation summary
- Quick start guide

### `QUICK_START.md`
- Step-by-step setup
- Manual testing examples
- Troubleshooting
- Pro tips

### `API_DOCUMENTATION.md`
- Complete API reference
- Request/response examples
- Error codes
- Signature chain explanation

### `IMPLEMENTATION.md`
- Architecture decisions
- Design patterns
- Concurrency strategy
- Future enhancements

### `PROJECT_STRUCTURE.md`
- This file
- Package overview
- Data flow diagrams
- Dependency graph

---

## Configuration

### `go.mod`
```go
module github.com/fiskaly/coding-challenges/signing-service-challenge

go 1.20
```

**Dependencies**: Only standard library (no external dependencies!)

---

## Build Artifacts

### `signature-service` (binary)
- Compiled application
- Ready to run
- ~10MB size

---

## Code Statistics

```
Package         Files    Lines    Tests
─────────────────────────────────────────
api/            4        ~400     9 tests
crypto/         4        ~200     -
domain/         4        ~350     5 tests
persistence/    1        ~90      -
main.go         1        ~25      -
─────────────────────────────────────────
Total           14       ~1,065   14 tests
Documentation   6        ~2,000   -
```

---

## Learning Path

**For understanding the codebase**:

1. Start with `README.md` (overview)
2. Read `QUICK_START.md` (hands-on)
3. Explore `main.go` (entry point)
4. Study `domain/device.go` (core entity)
5. Review `domain/service.go` (business logic)
6. Check `api/device.go` (HTTP layer)
7. Examine `persistence/inmemory.go` (storage)
8. Understand `crypto/signer.go` (algorithms)
9. Read tests for usage examples
10. Deep dive with `IMPLEMENTATION.md`

---

## Extension Points

### Adding a New Algorithm
1. Implement `crypto.Signer` interface
2. Add constant in `domain/device.go`
3. Update factory in `domain/service.go`

### Adding Database Support
1. Implement `domain.DeviceRepository` interface
2. Use marshalers in `crypto/` package
3. Update `main.go` dependency injection

### Adding Authentication
1. Create middleware in `api/server.go`
2. Add JWT validation
3. Extract user context

### Adding Audit Logging
1. Create `audit/` package
2. Implement observer pattern
3. Log in `domain/service.go`

---

## Best Practices Demonstrated

- ✅ Clean Architecture
- ✅ SOLID Principles
- ✅ Repository Pattern
- ✅ Strategy Pattern
- ✅ Dependency Injection
- ✅ Interface-based Design
- ✅ Comprehensive Testing
- ✅ Clear Documentation
- ✅ Thread Safety
- ✅ Error Handling

---

**Navigation**: See `README.md` for main documentation
