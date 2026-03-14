package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/octguy/auth-sqlc/internal/middleware"
	"github.com/octguy/auth-sqlc/internal/model"
	"github.com/octguy/auth-sqlc/internal/service"
	"github.com/octguy/auth-sqlc/pkg/response"
)

// AuthHandler holds the HTTP handlers for all auth routes.
type AuthHandler struct {
	authSvc service.AuthService
}

// NewAuthHandler constructs the handler with its service injected.
func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "registration failed")
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", resp)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "login failed")
		return
	}

	response.Success(c, http.StatusOK, "login successful", resp)
}

// GetProfile handles GET /api/v1/auth/profile (protected)
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.authSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	response.Success(c, http.StatusOK, "profile retrieved", user)
}
