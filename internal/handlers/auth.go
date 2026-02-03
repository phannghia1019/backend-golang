package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hbc/app-hrm-api-go/internal/middleware"
	"github.com/hbc/app-hrm-api-go/internal/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	permissionService *services.PermissionService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		permissionService: services.NewPermissionService(),
	}
}

// GetPermissions returns user permissions
// GET /v1/auth/permissions
func (h *AuthHandler) GetPermissions(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}

	permissions := h.permissionService.GetUserPermissions(user.Permissions)
	c.JSON(http.StatusOK, permissions)
}
