package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func setupTestServer() *Server {
	repository := persistence.NewInMemoryDeviceRepository()
	deviceService := domain.NewDeviceService(repository)
	return NewServer(":8080", deviceService)
}

func TestCreateDevice_Success(t *testing.T) {
	server := setupTestServer()

	reqBody := map[string]interface{}{
		"id":        "device-001",
		"algorithm": "RSA",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.CreateDevice(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	deviceData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a map")
	}

	if deviceData["id"] != "device-001" {
		t.Errorf("Expected device ID 'device-001', got %v", deviceData["id"])
	}

	if deviceData["algorithm"] != "RSA" {
		t.Errorf("Expected algorithm 'RSA', got %v", deviceData["algorithm"])
	}

	if deviceData["signature_counter"] != float64(0) {
		t.Errorf("Expected signature_counter 0, got %v", deviceData["signature_counter"])
	}
}

func TestCreateDevice_InvalidAlgorithm(t *testing.T) {
	server := setupTestServer()

	reqBody := map[string]interface{}{
		"id":        "device-002",
		"algorithm": "INVALID",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.CreateDevice(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateDevice_DuplicateID(t *testing.T) {
	server := setupTestServer()

	reqBody := map[string]interface{}{
		"id":        "device-003",
		"algorithm": "ECC",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(reqBody)

	// Create first device
	req1 := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	w1 := httptest.NewRecorder()
	server.CreateDevice(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("First creation failed with status %d", w1.Code)
	}

	// Try to create duplicate
	req2 := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	w2 := httptest.NewRecorder()
	server.CreateDevice(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate, got %d", http.StatusConflict, w2.Code)
	}
}

func TestGetDevice_Success(t *testing.T) {
	server := setupTestServer()

	// Create a device first
	reqBody := map[string]interface{}{
		"id":        "device-004",
		"algorithm": "RSA",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(reqBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	server.CreateDevice(createW, createReq)

	// Get the device
	getReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices/device-004", nil)
	getW := httptest.NewRecorder()
	server.GetDevice(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, getW.Code)
	}

	var response Response
	if err := json.Unmarshal(getW.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	deviceData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a map")
	}

	if deviceData["id"] != "device-004" {
		t.Errorf("Expected device ID 'device-004', got %v", deviceData["id"])
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/api/v0/devices/nonexistent", nil)
	w := httptest.NewRecorder()
	server.GetDevice(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListDevices(t *testing.T) {
	server := setupTestServer()

	// Create multiple devices
	devices := []string{"device-005", "device-006", "device-007"}
	for _, deviceID := range devices {
		reqBody := map[string]interface{}{
			"id":        deviceID,
			"algorithm": "RSA",
			"label":     "Test Device",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.CreateDevice(w, req)
	}

	// List devices
	listReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices", nil)
	listW := httptest.NewRecorder()
	server.ListDevices(listW, listReq)

	if listW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, listW.Code)
	}

	var response Response
	if err := json.Unmarshal(listW.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	deviceList, ok := response.Data.([]interface{})
	if !ok {
		t.Fatal("Response data is not an array")
	}

	if len(deviceList) != 3 {
		t.Errorf("Expected 3 devices, got %d", len(deviceList))
	}
}

func TestSignTransaction_Success(t *testing.T) {
	server := setupTestServer()

	// Create a device
	createBody := map[string]interface{}{
		"id":        "device-008",
		"algorithm": "ECC",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	server.CreateDevice(createW, createReq)

	// Sign transaction
	signBody := map[string]interface{}{
		"data": "transaction-data-123",
	}
	signBodyBytes, _ := json.Marshal(signBody)
	signReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices/device-008/sign", bytes.NewReader(signBodyBytes))
	signW := httptest.NewRecorder()
	server.SignTransaction(signW, signReq)

	if signW.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, signW.Code, signW.Body.String())
	}

	var response Response
	if err := json.Unmarshal(signW.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	signatureData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a map")
	}

	if signatureData["signature"] == "" {
		t.Error("Signature should not be empty")
	}

	if signatureData["signed_data"] == "" {
		t.Error("Signed data should not be empty")
	}
}

func TestSignTransaction_MonotonicCounter(t *testing.T) {
	server := setupTestServer()

	// Create a device
	createBody := map[string]interface{}{
		"id":        "device-009",
		"algorithm": "RSA",
		"label":     "Test Device",
	}
	body, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	server.CreateDevice(createW, createReq)

	// Perform multiple signatures and verify counter increments
	for i := 0; i < 5; i++ {
		signBody := map[string]interface{}{
			"data": "transaction-" + string(rune(i)),
		}
		signBodyBytes, _ := json.Marshal(signBody)
		signReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices/device-009/sign", bytes.NewReader(signBodyBytes))
		signW := httptest.NewRecorder()
		server.SignTransaction(signW, signReq)

		if signW.Code != http.StatusOK {
			t.Fatalf("Sign %d failed with status %d", i, signW.Code)
		}
	}

	// Get device and verify counter
	getReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices/device-009", nil)
	getW := httptest.NewRecorder()
	server.GetDevice(getW, getReq)

	var response Response
	json.Unmarshal(getW.Body.Bytes(), &response)
	deviceData := response.Data.(map[string]interface{})

	if deviceData["signature_counter"] != float64(5) {
		t.Errorf("Expected signature_counter 5, got %v", deviceData["signature_counter"])
	}
}

func TestSignTransaction_ConcurrentAccess(t *testing.T) {
	server := setupTestServer()

	// Create a device
	createBody := map[string]interface{}{
		"id":        "device-010",
		"algorithm": "RSA",
		"label":     "Concurrent Test Device",
	}
	body, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	server.CreateDevice(createW, createReq)

	// Perform concurrent signatures
	numGoroutines := 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			signBody := map[string]interface{}{
				"data": "concurrent-transaction",
			}
			signBodyBytes, _ := json.Marshal(signBody)
			signReq := httptest.NewRequest(http.MethodPost, "/api/v0/devices/device-010/sign", bytes.NewReader(signBodyBytes))
			signW := httptest.NewRecorder()
			server.SignTransaction(signW, signReq)

			if signW.Code != http.StatusOK {
				t.Errorf("Concurrent sign %d failed with status %d", index, signW.Code)
			}
		}(i)
	}

	wg.Wait()

	// Verify final counter is exactly numGoroutines (no gaps, no duplicates)
	getReq := httptest.NewRequest(http.MethodGet, "/api/v0/devices/device-010", nil)
	getW := httptest.NewRecorder()
	server.GetDevice(getW, getReq)

	var response Response
	json.Unmarshal(getW.Body.Bytes(), &response)
	deviceData := response.Data.(map[string]interface{})

	if deviceData["signature_counter"] != float64(numGoroutines) {
		t.Errorf("Expected signature_counter %d, got %v", numGoroutines, deviceData["signature_counter"])
	}
}
