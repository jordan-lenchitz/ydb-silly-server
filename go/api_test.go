package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVMStatusEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/vm/status", nil)
	req.Header.Set("X-VM-ID", "test-vm-123")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	var response VMState
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status == "" {
		t.Error("Expected Status in response, got empty string")
	}

	if response.System == nil {
		t.Error("Expected System map in response, got nil")
	} else {
		if _, ok := response.System["engine"]; !ok {
			t.Error("Expected 'engine' field in system metadata")
		}
		if _, ok := response.System["uptime"]; !ok {
			t.Error("Expected 'uptime' field in system metadata")
		}
	}
}

func TestExecuteNilSafety(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupRouter()

	req, _ := http.NewRequest("POST", "/api/execute", nil)
	w := httptest.NewRecorder()
	
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %v", w.Code)
	}
}