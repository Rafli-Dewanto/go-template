package handler

import (
	"encoding/json"
	"net/http"

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

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	writeResponse(w, statusCode, nil, message, nil)
}

func writeValidationErrorResponse(w http.ResponseWriter, errors []utils.ValidationError) {
	writeResponse(w, http.StatusBadRequest, errors, "Validation failed", nil)
}
