package handler

import (
	"net/http"

	"github.com/go4s/iam/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONErr(c, http.StatusBadRequest, CodeBadRequest, err.Error(), nil)
		return
	}

	token, userInfo, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		JSONErr(c, http.StatusUnauthorized, CodeUnauthorized, err.Error(), nil)
		return
	}

	JSONOK(c, map[string]any{
		"token": token,
		"user":  userInfo,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		JSONErr(c, http.StatusUnauthorized, CodeUnauthorized, "user not found in context", nil)
		return
	}

	userInfo, err := h.authService.Me(username.(string))
	if err != nil {
		JSONErr(c, http.StatusInternalServerError, CodeInternalError, err.Error(), nil)
		return
	}

	JSONOK(c, userInfo)
}
