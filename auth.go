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

// AddBearerAuth creates a Bearer token authentication hook
func AddBearerAuth(token string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// AddAPIKeyAuth creates an API key authentication hook
func AddAPIKeyAuth(apiKey, headerName string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set(headerName, apiKey)
		return nil
	}
}

// AddBasicAuth creates a Basic authentication hook
func AddBasicAuth(username, password string) RequestHook {
	return func(req *http.Request) error {
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		req.Header.Set("Authorization", "Basic "+auth)
		return nil
	}
}

// AddCustomAuth creates a custom authentication middleware
func AddCustomAuth(headerName, headerValue string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set(headerName, headerValue)
		return nil
	}
}

// AddAuthentication adds authentication headers using an AuthProvider
// For simple cases, use AddBearerAuth, AddAPIKeyAuth, etc. directly
func AddAuthentication(authProvider AuthProvider) RequestHook {
	return func(req *http.Request) error {
		token, err := authProvider.GetToken(req.Context())
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if token != "" {
			authHeader := authProvider.GetAuthHeader(token)
			req.Header.Set("Authorization", authHeader)
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
		return AddCustomAuth(config.HeaderName, config.HeaderValue)
	default:
		return func(req *http.Request) error {
			return nil
		}
	}
}

// AddDynamicAuth creates a dynamic authentication hook
func AddDynamicAuth(authFunc func(*http.Request) error) RequestHook {
	return func(req *http.Request) error {
		return authFunc(req)
	}
}

// AddJWTTokenAuth creates a JWT token authentication middleware
func AddJWTTokenAuth(token string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// AddOAuth2Auth creates an OAuth2 authentication middleware
func AddOAuth2Auth(accessToken string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}
}

// AddAPIKeyHeaderAuth creates an API key authentication middleware with common header names
func AddAPIKeyHeaderAuth(apiKey string) RequestHook {
	return func(req *http.Request) error {
		// Try common API key header names
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("API-Key", apiKey)
		return nil
	}
}

// AddDigestAuth creates a Digest authentication middleware (simplified)
func AddDigestAuth(username, password, realm string) RequestHook {
	return func(req *http.Request) error {
		// This is a simplified implementation
		// In a real implementation, you'd need to handle the digest challenge
		req.Header.Set("Authorization", "Digest username=\""+username+"\", realm=\""+realm+"\"")
		return nil
	}
}

// AddMultiAuth creates a middleware that tries multiple authentication methods
func AddMultiAuth(authMethods ...RequestHook) RequestHook {
	return func(req *http.Request) error {
		for _, authMethod := range authMethods {
			if err := authMethod(req); err != nil {
				return err
			}
		}
		return nil
	}
}

// AddConditionalAuth creates a conditional authentication middleware
func AddConditionalAuth(condition func(*http.Request) bool, authMiddleware RequestHook) RequestHook {
	return func(req *http.Request) error {
		if condition(req) {
			return authMiddleware(req)
		}
		return nil
	}
}

// AddAuthFromContext creates an authentication middleware that gets credentials from context
func AddAuthFromContext(headerName string, contextKey string) RequestHook {
	return func(req *http.Request) error {
		if token, ok := req.Context().Value(contextKey).(string); ok {
			req.Header.Set(headerName, token)
		}
		return nil
	}
}

// AddRotatingAuth creates a rotating authentication middleware
func AddRotatingAuth(tokens []string, headerName string) RequestHook {
	counter := 0
	return func(req *http.Request) error {
		if len(tokens) == 0 {
			return nil
		}

		token := tokens[counter%len(tokens)]
		req.Header.Set(headerName, token)
		counter++

		return nil
	}
}
