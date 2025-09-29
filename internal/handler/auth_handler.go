package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rafli-Dewanto/go-template/internal/auth"
	"github.com/Rafli-Dewanto/go-template/internal/context"
	"github.com/Rafli-Dewanto/go-template/internal/model"
	"github.com/Rafli-Dewanto/go-template/internal/service"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userService  service.UserService
	tokenManager *auth.TokenManager
	logger       *utils.Logger
}

func NewAuthHandler(userService service.UserService, tokenManager *auth.TokenManager, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{userService: userService, tokenManager: tokenManager, logger: logger}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	apiID := context.GetAPIID(ctx)
	cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		h.logger.WarningWithAPIID(apiID, "Method not allowed: %s", r.Method)
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorWithAPIID(apiID, "Failed to decode request body: %v", err)
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.GetByEmail(cancelCtx, req.Email)
	if err != nil {
		h.logger.WarningWithAPIID(apiID, "Invalid credentials")
		WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := auth.ComparePassword(user.Password, req.Password); err != nil {
		h.logger.WarningWithAPIID(apiID, "Invalid credentials")
		WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.tokenManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		h.logger.ErrorWithAPIID(apiID, "Failed to generate token: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := LoginResponse{Token: token}
	writeResponse(w, http.StatusOK, response, "Login successful", nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithRequestID(r.Context(), uuid.New().String())
	apiID := context.GetAPIID(ctx)
	cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if r.Method != http.MethodPost {
		h.logger.WarningWithAPIID(apiID, "Method not allowed: %s", r.Method)
		WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorWithAPIID(apiID, "Failed to decode request body: %v", err)
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.ErrorWithAPIID(apiID, "Failed to hash password: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	req.Password = hashedPassword

	err = h.userService.Create(cancelCtx, &req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			h.logger.WarningWithAPIID(apiID, "Invalid input for user registration: %v", err)
			WriteErrorResponse(w, http.StatusBadRequest, "Invalid input")
			return
		case service.ErrUserAlreadyExists:
			h.logger.WarningWithAPIID(apiID, "User with email or username already exists")
			WriteErrorResponse(w, http.StatusConflict, "User already exists")
			return
		default:
			h.logger.ErrorWithAPIID(apiID, "Failed to create user: %v", err)
			WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	writeResponse(w, http.StatusCreated, nil, "User registered successfully", nil)
}
