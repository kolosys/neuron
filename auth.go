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

// AddBearerAuth creates a Bearer token authentication middleware
func AddBearerAuth(token string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// AddAPIKeyAuth creates an API key authentication middleware
func AddAPIKeyAuth(apiKey, headerName string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set(headerName, apiKey)
		return nil
	}
}

// AddBasicAuth creates a Basic authentication middleware
func AddBasicAuth(username, password string) RequestMiddleware {
	return func(req *http.Request) error {
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		req.Header.Set("Authorization", "Basic "+auth)
		return nil
	}
}

// AddCustomAuth creates a custom authentication middleware
func AddCustomAuth(headerName, headerValue string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set(headerName, headerValue)
		return nil
	}
}

// AddAuth creates an authentication middleware based on configuration
func AddAuth(config AuthConfig) RequestMiddleware {
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

// AddDynamicAuth creates a dynamic authentication middleware
func AddDynamicAuth(authFunc func(*http.Request) error) RequestMiddleware {
	return func(req *http.Request) error {
		return authFunc(req)
	}
}

// AddJWTTokenAuth creates a JWT token authentication middleware
func AddJWTTokenAuth(token string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// AddOAuth2Auth creates an OAuth2 authentication middleware
func AddOAuth2Auth(accessToken string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}
}

// AddAPIKeyHeaderAuth creates an API key authentication middleware with common header names
func AddAPIKeyHeaderAuth(apiKey string) RequestMiddleware {
	return func(req *http.Request) error {
		// Try common API key header names
		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("API-Key", apiKey)
		return nil
	}
}

// AddDigestAuth creates a Digest authentication middleware (simplified)
func AddDigestAuth(username, password, realm string) RequestMiddleware {
	return func(req *http.Request) error {
		// This is a simplified implementation
		// In a real implementation, you'd need to handle the digest challenge
		req.Header.Set("Authorization", "Digest username=\""+username+"\", realm=\""+realm+"\"")
		return nil
	}
}

// AddMultiAuth creates a middleware that tries multiple authentication methods
func AddMultiAuth(authMethods ...RequestMiddleware) RequestMiddleware {
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
func AddConditionalAuth(condition func(*http.Request) bool, authMiddleware RequestMiddleware) RequestMiddleware {
	return func(req *http.Request) error {
		if condition(req) {
			return authMiddleware(req)
		}
		return nil
	}
}

// AddAuthFromContext creates an authentication middleware that gets credentials from context
func AddAuthFromContext(headerName string, contextKey string) RequestMiddleware {
	return func(req *http.Request) error {
		if token, ok := req.Context().Value(contextKey).(string); ok {
			req.Header.Set(headerName, token)
		}
		return nil
	}
}

// AddRotatingAuth creates a rotating authentication middleware
func AddRotatingAuth(tokens []string, headerName string) RequestMiddleware {
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
