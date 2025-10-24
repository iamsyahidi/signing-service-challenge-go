package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

// CreateDevice handles POST /api/v0/devices
func (s *Server) CreateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Failed to read request body"})
		return
	}
	defer r.Body.Close()

	var req domain.CreateSignatureDeviceRequest
	if err := json.Unmarshal(body, &req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON format"})
		return
	}

	// Create device
	device, err := s.deviceService.CreateDevice(&req)
	if err != nil {
		if errors.Is(err, persistence.ErrDeviceAlreadyExists) {
			WriteErrorResponse(w, http.StatusConflict, []string{err.Error()})
			return
		}
		WriteErrorResponse(w, http.StatusBadRequest, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusCreated, device)
}

// GetDevice handles GET /api/v0/devices/{id}
func (s *Server) GetDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	// Extract device ID from path
	deviceID := strings.TrimPrefix(r.URL.Path, "/api/v0/devices/")
	if deviceID == "" || deviceID == r.URL.Path {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Device ID is required"})
		return
	}

	// Get device
	device, err := s.deviceService.GetDevice(deviceID)
	if err != nil {
		if errors.Is(err, persistence.ErrDeviceNotFound) {
			WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})
			return
		}
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusOK, device)
}

// ListDevices handles GET /api/v0/devices
func (s *Server) ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	devices, err := s.deviceService.ListDevices()
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusOK, devices)
}

// SignTransaction handles POST /api/v0/devices/{id}/sign
func (s *Server) SignTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	// Extract device ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v0/devices/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "sign" {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid path format"})
		return
	}
	deviceID := parts[0]

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Failed to read request body"})
		return
	}
	defer r.Body.Close()

	var req domain.SignTransactionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, []string{"Invalid JSON format"})
		return
	}

	// Sign transaction
	response, err := s.deviceService.SignTransaction(deviceID, req.Data)
	if err != nil {
		if errors.Is(err, persistence.ErrDeviceNotFound) {
			WriteErrorResponse(w, http.StatusNotFound, []string{err.Error()})
			return
		}
		WriteErrorResponse(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	WriteAPIResponse(w, http.StatusOK, response)
}

// DevicesRouter routes device-related requests
func (s *Server) DevicesRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// POST /api/v0/devices - Create device
	if path == "/api/v0/devices" && r.Method == http.MethodPost {
		s.CreateDevice(w, r)
		return
	}

	// GET /api/v0/devices - List devices
	if path == "/api/v0/devices" && r.Method == http.MethodGet {
		s.ListDevices(w, r)
		return
	}

	// POST /api/v0/devices/{id}/sign - Sign transaction
	if strings.HasSuffix(path, "/sign") && r.Method == http.MethodPost {
		s.SignTransaction(w, r)
		return
	}

	// GET /api/v0/devices/{id} - Get device
	if strings.HasPrefix(path, "/api/v0/devices/") && r.Method == http.MethodGet {
		deviceID := strings.TrimPrefix(path, "/api/v0/devices/")
		if !strings.Contains(deviceID, "/") {
			s.GetDevice(w, r)
			return
		}
	}

	WriteErrorResponse(w, http.StatusNotFound, []string{
		fmt.Sprintf("Route not found: %s %s", r.Method, path),
	})
}
