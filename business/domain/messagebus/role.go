package messagebus

import "fmt"

// Set of known userRoles.
var userRoles = make(map[string]Role)

// Set of known llmRoles.
var llmRoles = make(map[string]Role)

// Set of possible roles for a user.
var (
	RoleUser = newUserRole("user")
)

var (
	RoleSystem    = newLlmRole("system")
	RoleAssistant = newLlmRole("assistant")
)

// Role represents a role in the system.
type Role struct {
	name string
}

func NewRole(name string) Role {
	return Role{name}
}

func newUserRole(role string) Role {
	r := Role{role}
	userRoles[role] = r
	return r
}

func newLlmRole(role string) Role {
	r := Role{role}
	llmRoles[role] = r
	return r
}

// ParseUserRoles parses the string value and returns a role if one exists.
func ParseUserRoles(value string) (Role, error) {
	role, exists := userRoles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

// ParseUserRoles parses the string value and returns a role if one exists.
func ParseLlmRoles(value string) (Role, error) {
	role, exists := llmRoles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

// Name returns the name of the role.
func (r Role) Name() string {
	return r.name
}
