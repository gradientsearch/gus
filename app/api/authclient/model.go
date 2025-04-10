package authclient

import (
	"github.com/google/uuid"
	"github.com/gradientsearch/gus/app/api/auth"
)

// Error represents an error in the system.
type Error struct {
	Message string `json:"message"`
}

// Error implements the error interface.
func (err Error) Error() string {
	return err.Message
}

// Authorize defines the information required to perform an authorization.
type Authorize struct {
	Claims auth.Claims
	UserID uuid.UUID
	Rule   string
}

// AuthenticateResp defines the information that will be received on authenticate.
type AuthenticateResp struct {
	UserID uuid.UUID
	Claims auth.Claims
}
