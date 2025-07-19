package auth

import (
	"errors"
	"net/http"
)

// Define auth errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoCredentials      = errors.New("no credentials provided")
)

// User represents an authenticated user
type AuthUser struct {
	ID       string
	Username string
}

// Authenticator is the interface that all auth methods must implement
type Authenticator interface {
	// Authenticate checks the request for valid credentials or tokens
	Authenticate(r *http.Request) (*AuthUser, error)
}
