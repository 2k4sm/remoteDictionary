package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2k4sm/remoteDictionary/src/cache"
	"github.com/2k4sm/remoteDictionary/src/models"
)

func setupHandler() *Handler {
	// Initialize cache with small limits for testing
	c := cache.NewCache(100, 100)
	return NewHandler(c)
}

func TestPutHandler(t *testing.T) {
	h := setupHandler()

	tests := []struct {
		name           string
		method         string
		body           models.PutRequest
		expectedStatus int
	}{
		{
			name:           "Valid Put",
			method:         http.MethodPost,
			body:           models.PutRequest{Key: "test", Value: "value"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodGet,
			body:           models.PutRequest{Key: "test", Value: "value"},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Empty Key",
			method:         http.MethodPost,
			body:           models.PutRequest{Key: "", Value: "value"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/put", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.PutHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	h := setupHandler()
	// Pre-populate cache
	_ = h.cache.Put("existing", "data")

	tests := []struct {
		name           string
		method         string
		key            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Get",
			method:         http.MethodGet,
			key:            "existing",
			expectedStatus: http.StatusOK,
			expectedBody:   "data",
		},
		{
			name:           "Key Not Found",
			method:         http.MethodGet,
			key:            "missing",
			expectedStatus: http.StatusOK, // Handler returns 200 with ERROR status in body
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			key:            "existing",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Missing Key Parameter",
			method:         http.MethodGet,
			key:            "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/get?key="+tt.key, nil)
			w := httptest.NewRecorder()

			h.GetHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.expectedBody != "" && w.Code == http.StatusOK {
				var resp models.GetResponse
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Value != tt.expectedBody {
					t.Errorf("expected body %s, got %s", tt.expectedBody, resp.Value)
				}
			}
		})
	}
}
