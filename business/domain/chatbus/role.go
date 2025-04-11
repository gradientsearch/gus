package chatbus

import "fmt"

// Set of known roles.
var roles = make(map[string]Role)

// Set of possible roles for a user.
var (
	RoleUser = newRole("user")
)

var RoleAssistant = Role{"assistant"}

// Role represents a role in the system.
type Role struct {
	name string
}

func newRole(role string) Role {
	r := Role{role}
	roles[role] = r
	return r
}

// ParseRole parses the string value and returns a role if one exists.
func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

// Name returns the name of the role.
func (r Role) Name() string {
	return r.name
}
