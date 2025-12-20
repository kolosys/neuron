package neuron

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

// AuthConfig configures authentication hooks
type AuthConfig struct {
	Type        AuthType
	Token       string
	Username    string
	Password    string
	HeaderName  string
	HeaderValue string
}

// AuthType represents the type of authentication
type AuthType int

const (
	AuthTypeBearer AuthType = iota
	AuthTypeAPIKey
	AuthTypeBasic
	AuthTypeCustom
)

// AddAuthHeader creates a generic header-based authentication hook
func AddAuthHeader(header, value string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set(header, value)
		return nil
	}
}

// AddBearerAuth creates a Bearer token authentication hook
func AddBearerAuth(token string) RequestHook {
	return AddAuthHeader("Authorization", "Bearer "+token)
}

// AddAPIKeyAuth creates an API key authentication hook
func AddAPIKeyAuth(apiKey, headerName string) RequestHook {
	return AddAuthHeader(headerName, apiKey)
}

// AddBasicAuth creates a Basic authentication hook
func AddBasicAuth(username, password string) RequestHook {
	return func(req *http.Request) error {
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		req.Header.Set("Authorization", "Basic "+auth)
		return nil
	}
}

// AddAuthentication adds authentication headers using an AuthProvider
func AddAuthentication(authProvider AuthProvider) RequestHook {
	return func(req *http.Request) error {
		token, err := authProvider.GetToken(req.Context())
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if token != "" {
			req.Header.Set("Authorization", authProvider.GetAuthHeader(token))
		}

		return nil
	}
}

// AuthProvider interface for authentication
type AuthProvider interface {
	GetToken(ctx context.Context) (string, error)
	GetAuthHeader(token string) string
}

// StaticAuthProvider provides static token authentication
type StaticAuthProvider struct {
	Token  string
	Prefix string
}

func (a *StaticAuthProvider) GetToken(ctx context.Context) (string, error) {
	return a.Token, nil
}

func (a *StaticAuthProvider) GetAuthHeader(token string) string {
	if a.Prefix != "" {
		return a.Prefix + " " + token
	}
	return "Bearer " + token
}

// AddAuth creates an authentication hook based on configuration
func AddAuth(config AuthConfig) RequestHook {
	switch config.Type {
	case AuthTypeBearer:
		return AddBearerAuth(config.Token)
	case AuthTypeAPIKey:
		return AddAPIKeyAuth(config.Token, config.HeaderName)
	case AuthTypeBasic:
		return AddBasicAuth(config.Username, config.Password)
	case AuthTypeCustom:
		return AddAuthHeader(config.HeaderName, config.HeaderValue)
	default:
		return func(req *http.Request) error { return nil }
	}
}

// AddAPIKeyHeaderAuth creates an API key authentication middleware with common header names
func AddAPIKeyHeaderAuth(apiKey string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("API-Key", apiKey)
		return nil
	}
}
