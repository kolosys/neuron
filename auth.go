package neuron

import (
	"encoding/base64"
	"net/http"
)

// AuthConfig configures authentication middleware
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

// BearerAuthMiddleware creates a Bearer token authentication middleware
func BearerAuthMiddleware(token string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// APIKeyAuthMiddleware creates an API key authentication middleware
func APIKeyAuthMiddleware(apiKey, headerName string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set(headerName, apiKey)
		return nil
	}
}

// BasicAuthMiddleware creates a Basic authentication middleware
func BasicAuthMiddleware(username, password string) RequestMiddleware {
	return func(req *http.Request) error {
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		req.Header.Set("Authorization", "Basic "+auth)
		return nil
	}
}

// CustomAuthMiddleware creates a custom authentication middleware
func CustomAuthMiddleware(headerName, headerValue string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set(headerName, headerValue)
		return nil
	}
}

// AuthMiddleware creates an authentication middleware based on configuration
func AuthMiddleware(config AuthConfig) RequestMiddleware {
	switch config.Type {
	case AuthTypeBearer:
		return BearerAuthMiddleware(config.Token)
	case AuthTypeAPIKey:
		return APIKeyAuthMiddleware(config.Token, config.HeaderName)
	case AuthTypeBasic:
		return BasicAuthMiddleware(config.Username, config.Password)
	case AuthTypeCustom:
		return CustomAuthMiddleware(config.HeaderName, config.HeaderValue)
	default:
		return func(req *http.Request) error {
			return nil
		}
	}
}

// DynamicAuthMiddleware creates a dynamic authentication middleware
func DynamicAuthMiddleware(authFunc func(*http.Request) error) RequestMiddleware {
	return func(req *http.Request) error {
		return authFunc(req)
	}
}

// JWTTokenAuthMiddleware creates a JWT token authentication middleware
func JWTTokenAuthMiddleware(token string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// OAuth2AuthMiddleware creates an OAuth2 authentication middleware
func OAuth2AuthMiddleware(accessToken string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}
}

// APIKeyHeaderAuthMiddleware creates an API key authentication middleware with common header names
func APIKeyHeaderAuthMiddleware(apiKey string) RequestMiddleware {
	return func(req *http.Request) error {
		// Try common API key header names
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("API-Key", apiKey)
		return nil
	}
}

// DigestAuthMiddleware creates a Digest authentication middleware (simplified)
func DigestAuthMiddleware(username, password, realm string) RequestMiddleware {
	return func(req *http.Request) error {
		// This is a simplified implementation
		// In a real implementation, you'd need to handle the digest challenge
		req.Header.Set("Authorization", "Digest username=\""+username+"\", realm=\""+realm+"\"")
		return nil
	}
}

// MultiAuthMiddleware creates a middleware that tries multiple authentication methods
func MultiAuthMiddleware(authMethods ...RequestMiddleware) RequestMiddleware {
	return func(req *http.Request) error {
		for _, authMethod := range authMethods {
			if err := authMethod(req); err != nil {
				return err
			}
		}
		return nil
	}
}

// ConditionalAuthMiddleware creates a conditional authentication middleware
func ConditionalAuthMiddleware(condition func(*http.Request) bool, authMiddleware RequestMiddleware) RequestMiddleware {
	return func(req *http.Request) error {
		if condition(req) {
			return authMiddleware(req)
		}
		return nil
	}
}

// AuthFromContextMiddleware creates an authentication middleware that gets credentials from context
func AuthFromContextMiddleware(headerName string, contextKey string) RequestMiddleware {
	return func(req *http.Request) error {
		if token, ok := req.Context().Value(contextKey).(string); ok {
			req.Header.Set(headerName, token)
		}
		return nil
	}
}

// RotatingAuthMiddleware creates a rotating authentication middleware
func RotatingAuthMiddleware(tokens []string, headerName string) RequestMiddleware {
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
