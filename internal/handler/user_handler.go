package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Rafli-Dewanto/go-template/internal/context"
	"github.com/Rafli-Dewanto/go-template/internal/model"
	"github.com/Rafli-Dewanto/go-template/internal/service"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	logger      *utils.Logger
}

func NewUserHandler(userService service.UserService, logger *utils.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		h.logger.Warning("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateUserRequest // Use value instead of pointer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.userService.Create(cancelCtx, &req)
	if err != nil {
		// Properly check for context deadline exceeded error
		if errors.Is(err, ctx.Err()) {
			h.logger.Warning("Request timeout: operation took longer than 10 seconds")
			writeErrorResponse(w, http.StatusRequestTimeout, "Request timeout")
			return
		}

		switch err {
		case service.ErrInvalidInput:
			h.logger.Warning("Invalid input for user creation: %v", err)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid input")
		case service.ErrUserAlreadyExists:
			h.logger.Warning("User with email or username already exists")
			writeErrorResponse(w, http.StatusConflict, "User already exists")
		default:
			h.logger.Error("Failed to create user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	writeResponse(w, http.StatusCreated, nil, "User created successfully", nil)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	if r.Method != http.MethodGet {
		h.logger.Warning("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" || strings.Contains(idStr, "/") {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Failed to parse user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrUserNotFound {
			h.logger.Warning("User not found with ID: %d", id)
			writeErrorResponse(w, http.StatusNotFound, "User not found")
			return
		}
		h.logger.Error("Failed to get user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	writeResponse(w, http.StatusOK, user, "User retrieved successfully", nil)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	if r.Method != http.MethodGet {
		h.logger.Warning("Method not allowed: %s", r.Method)
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	h.logger.Info("List users with limit: %d, offset: %d", limit, offset)

	query := &model.PaginationQuery{
		Page:   utils.Default(page, 1),
		Limit:  utils.Default(limit, 10),
		Offset: utils.Default(offset, 0),
	}

	response, err := h.userService.List(ctx, query)
	if err != nil {
		h.logger.Error("Failed to list users: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, response.Data, response.Message, response.Meta)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	if r.Method != http.MethodPut {
		h.logger.Warning("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" || strings.Contains(idStr, "/") {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Failed to parse user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req *model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.ID = id

	err = h.userService.Update(ctx, req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			h.logger.Warning("Invalid input for user update: %v", err)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid input")
		case service.ErrUserNotFound:
			h.logger.Warning("User not found for update with ID: %d", req.ID)
			writeErrorResponse(w, http.StatusNotFound, "User not found")
		default:
			h.logger.Error("Failed to update user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	writeResponse(w, http.StatusOK, nil, "User updated successfully", nil)
}

func (h *UserHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	h.logger.Info("Handling delete user request")
	if r.Method != http.MethodPatch {
		h.logger.Warning("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Failed to parse user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.SoftDelete(ctx, id); err != nil {
		switch err {
		case service.ErrUserNotFound:
			h.logger.Warning("User not found for deletion with ID: %d", id)
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			h.logger.Error("Failed to delete user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
