package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/2k4sm/remoteDictionary/src/cache"
	"github.com/2k4sm/remoteDictionary/src/models"
)

var buffersPool = sync.Pool{
	New: func() any {
		b := make([]byte, 1024)
		return &b
	},
}

type Handler struct {
	cache *cache.Cache
}

func NewHandler(cache *cache.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bufPtr := buffersPool.Get().(*[]byte)
	defer buffersPool.Put(bufPtr)
	buf := *bufPtr

	body, err := io.ReadAll(io.LimitReader(r.Body, int64(len(buf))))
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		sendErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req models.PutRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Key) == "" {
		sendErrorResponse(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.cache.Put(req.Key, req.Value)
	if err != nil {
		switch err {
		case cache.ErrKeyTooLarge:
			sendErrorResponse(w, "Key exceeds maximum length (256 characters)", http.StatusBadRequest)
		case cache.ErrValueTooLarge:
			sendErrorResponse(w, "Value exceeds maximum length (256 characters)", http.StatusBadRequest)
		default:
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := models.PutResponse{
		Status:  "OK",
		Message: "Key inserted/updated successfully.",
	}
	sendJSONResponse(w, resp, http.StatusOK)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if strings.TrimSpace(key) == "" {
		sendErrorResponse(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	value, err := h.cache.Get(key)
	if err != nil {
		if err == cache.ErrKeyNotFound {
			resp := models.ErrorResponse{
				Status:  "ERROR",
				Message: "Key not found.",
			}
			sendJSONResponse(w, resp, http.StatusOK)
		} else {
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := models.GetResponse{
		Status:  "OK",
		Key:     key,
		Value:   value,
		Message: "Key retrieved successfully.",
	}
	sendJSONResponse(w, resp, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	resp := models.ErrorResponse{
		Status:  "ERROR",
		Message: message,
	}
	sendJSONResponse(w, resp, statusCode)
}
