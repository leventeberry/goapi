package services

// ValidRoles contains the list of allowed user roles
var ValidRoles = []string{"user", "admin"}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

