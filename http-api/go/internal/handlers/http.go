package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type httpError struct {
	Message string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", slog.Any("error", err))
	}
}
