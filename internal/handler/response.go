package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Rafli-Dewanto/go-template/internal/context"
	"github.com/Rafli-Dewanto/go-template/internal/model"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
)

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}, message string, meta *model.PaginatedMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if errors, ok := data.([]utils.ValidationError); ok {
			resp := model.FailedResponse{
				Message: message,
				Errors:  errors,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp := model.Response{
		Data:    data,
		Message: message,
		Meta:    meta,
	}

	json.NewEncoder(w).Encode(resp)
}

// writeResponseWithContext writes HTTP response with API ID from context
func writeResponseWithContext(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, message string, meta *model.PaginatedMeta) {
	// Get API ID from context and add to response header
	if apiID := context.GetAPIID(r.Context()); apiID != "" {
		w.Header().Set("X-API-ID", apiID)
	}

	writeResponse(w, statusCode, data, message, meta)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	writeResponse(w, statusCode, nil, message, nil)
}

// WriteErrorResponseWithContext writes error response with API ID from context
func WriteErrorResponseWithContext(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	writeResponseWithContext(w, r, statusCode, nil, message, nil)
}

func writeValidationErrorResponse(w http.ResponseWriter, errors []utils.ValidationError) {
	writeResponse(w, http.StatusBadRequest, errors, "Validation failed", nil)
}

// writeValidationErrorResponseWithContext writes validation error response with API ID from context
func writeValidationErrorResponseWithContext(w http.ResponseWriter, r *http.Request, errors []utils.ValidationError) {
	writeResponseWithContext(w, r, http.StatusBadRequest, errors, "Validation failed", nil)
}
