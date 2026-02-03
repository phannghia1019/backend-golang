package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hbc/app-hrm-api-go/internal/config"
	"github.com/hbc/app-hrm-api-go/internal/models"
	"github.com/hbc/app-hrm-api-go/internal/services"
	"github.com/hbc/app-hrm-api-go/internal/utils"
)

var authLog = utils.NewLogger("Auth")

// UserContext holds authenticated user information
type UserContext struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	EmailAddress  string   `json:"emailAddress"`
	Name          string   `json:"name"`
	CompanyDbName string   `json:"companyDbName"`
	CompanyName   string   `json:"companyName"`
	Permissions   []string `json:"permissions"`
	ManagerEmail  string   `json:"managerEmail"`
}

// ExternalAPIResponse represents response from App Manager API
type ExternalAPIResponse struct {
	EmailAddress string   `json:"emailAddress"`
	Permissions  []string `json:"permissions"`
	Name         string   `json:"name"`
}

// AuthMiddleware validates authentication and populates user context
func AuthMiddleware() gin.HandlerFunc {
	userLookup := services.NewUserLookupService()

	return func(c *gin.Context) {
		// Check authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			err := models.Unauthorized("authMiddleware", "Authorization header is required", "Missing Authorization header")
			models.HandleError(c, err)
			c.Abort()
			return
		}

		// Call external API to verify token
		cfg := config.Get()
		apiURL := fmt.Sprintf("%s/v1/user/authorize/permission", cfg.AppManagerAPIHost)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			authLog.Error("Failed to create request: %v", err)
			appErr := models.InternalError("authMiddleware", "Authentication system error", "Failed to create verification request")
			models.HandleError(c, appErr)
			c.Abort()
			return
		}

		req.Header.Set("Authorization", authHeader)
		q := req.URL.Query()
		q.Add("hostUrl", c.GetHeader("Origin"))
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			authLog.Error("Failed to call external API: %v", err)
			appErr := models.InternalError("authMiddleware", "Authentication service error", "Failed to verify token")
			models.HandleError(c, appErr)
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			authLog.Error("External API error (%d): %s", resp.StatusCode, string(body))
			appErr := models.NewAppError("authMiddleware", "Authentication failed", fmt.Sprintf("External API returned %d", resp.StatusCode), resp.StatusCode)
			models.HandleError(c, appErr)
			c.Abort()
			return
		}

		// Parse external API response
		var apiResponse ExternalAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			authLog.Error("Failed to parse API response: %v", err)
			appErr := models.InternalError("authMiddleware", "Authentication system error", "Invalid response from auth service")
			models.HandleError(c, appErr)
			c.Abort()
			return
		}

		if apiResponse.EmailAddress == "" {
			appErr := models.Forbidden("authMiddleware", "Email address not found in External API response", "External API must provide emailAddress field")
			models.HandleError(c, appErr)
			c.Abort()
			return
		}

		// Find user in SAP databases
		authLog.Info("Looking up user in SAP databases: %s", apiResponse.EmailAddress)

		userInfo, err := userLookup.FindUserAcrossAllDatabases(apiResponse.EmailAddress)
		if err != nil || userInfo == nil {
			authLog.Error("User not found in any SAP database: %s", apiResponse.EmailAddress)
			appErr := models.Forbidden("authMiddleware", "User not found in company databases",
				fmt.Sprintf("Email %s not found in any SAP database. Please contact administrator.", apiResponse.EmailAddress))
			models.HandleError(c, appErr)
			c.Abort()
			return
		}

		authLog.Info("Successfully found user in database: %s (%s)", userInfo.CompanyName, userInfo.CompanyDbName)
		authLog.Info("User details - Name: %s, Manager: %s", userInfo.Name, userInfo.ManagerEmail)

		// Validate permissions
		if apiResponse.Permissions == nil {
			apiResponse.Permissions = []string{}
		}

		if len(apiResponse.Permissions) == 0 {
			authLog.Warn("No permissions provided by External API for %s", apiResponse.EmailAddress)
		} else {
			authLog.Info("Permissions from External API: %s", strings.Join(apiResponse.Permissions, ", "))
		}

		// Set user context
		user := &UserContext{
			EmailAddress:  apiResponse.EmailAddress,
			Name:          userInfo.Name,
			CompanyDbName: userInfo.CompanyDbName,
			CompanyName:   userInfo.CompanyName,
			Permissions:   apiResponse.Permissions,
			ManagerEmail:  userInfo.ManagerEmail,
		}

		c.Set("user", user)
		authLog.Info("Authentication successful for: %s from %s", apiResponse.EmailAddress, userInfo.CompanyName)

		c.Next()
	}
}

// GetUser retrieves user from context
func GetUser(c *gin.Context) *UserContext {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*UserContext)
}
