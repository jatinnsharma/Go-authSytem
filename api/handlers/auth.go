package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jatinnsharma/internal/auth"
	"github.com/jatinnsharma/internal/utils"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req auth.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	user, err := h.authService.Signup(&req, userAgent, ipAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Signup failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	response, err := h.authService.Login(&req, userAgent, ipAddress)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Refresh token required")
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Refresh token required")
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
