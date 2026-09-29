// Package handler provides HTTP handlers for the application.
package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler struct
type HealthHandler struct{}

// NewHealtHandler constructor
func NewHealtHandler() *HealthHandler {
	return &HealthHandler{}
}

// healthResponse struct
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ServerHTTP function
func (handle *HealthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	resp := healthResponse{
		Status:  "ok",
		Version: "1.0.0",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK) // opcional

	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return
	}

}
