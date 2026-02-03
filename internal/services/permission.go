package services

// Available permissions in the system
const (
	// Wildcard permission (super admin)
	PermissionWildcard = "*"

	// Super admin - can access all companies
	PermissionSuperAdmin = "hbg_super_admin"
)

// AllPermissions contains all available permissions
var AllPermissions = []string{
	PermissionWildcard,
	PermissionSuperAdmin,
}

// PermissionService handles permission checks
type PermissionService struct{}

// NewPermissionService creates a new permission service
func NewPermissionService() *PermissionService {
	return &PermissionService{}
}

// HasPermission checks if user has specific permission
func (s *PermissionService) HasPermission(userPermissions []string, permission string) bool {
	if userPermissions == nil {
		return false
	}

	// Check wildcard permission first
	for _, p := range userPermissions {
		if p == PermissionWildcard {
			return true
		}
		if p == permission {
			return true
		}
	}
	return false
}

// IsSuperAdmin checks if user is super admin
func (s *PermissionService) IsSuperAdmin(userPermissions []string) bool {
	if userPermissions == nil {
		return false
	}

	for _, p := range userPermissions {
		if p == PermissionWildcard || p == PermissionSuperAdmin {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if user has any of the specified permissions
func (s *PermissionService) HasAnyPermission(userPermissions []string, permissions []string) bool {
	if userPermissions == nil {
		return false
	}

	// Check wildcard first
	for _, p := range userPermissions {
		if p == PermissionWildcard {
			return true
		}
	}

	// Check each permission
	for _, required := range permissions {
		for _, have := range userPermissions {
			if have == required {
				return true
			}
		}
	}
	return false
}

// HasWildcardPermission checks if user has wildcard permission
func (s *PermissionService) HasWildcardPermission(userPermissions []string) bool {
	for _, p := range userPermissions {
		if p == PermissionWildcard {
			return true
		}
	}
	return false
}

// GetUserPermissions returns user permissions, expanding wildcard if needed
func (s *PermissionService) GetUserPermissions(userPermissions []string) []string {
	if userPermissions == nil {
		return []string{}
	}

	// If user has wildcard, return all permissions
	for _, p := range userPermissions {
		if p == PermissionWildcard {
			return AllPermissions
		}
	}

	return userPermissions
}
