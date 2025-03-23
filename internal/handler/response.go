package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Rafli-Dewanto/go-template/internal/model"
)

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}, message string, meta *model.PaginatedMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := model.Response{
		Data:    data,
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Meta:    meta,
	}

	json.NewEncoder(w).Encode(resp)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	writeResponse(w, statusCode, nil, message, nil)
}
