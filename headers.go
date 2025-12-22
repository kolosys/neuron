package neuron

import (
	"net/http"
)

// HeaderConfig configures header hooks
type HeaderConfig struct {
	UserAgent     string
	CustomHeaders map[string]string
	RemoveHeaders []string
	AddHeaders    map[string]string
}

// DefaultHeaderConfig returns a default header configuration
func DefaultHeaderConfig() HeaderConfig {
	return HeaderConfig{
		UserAgent:     "Neuron-HTTP-Client/1.0",
		CustomHeaders: make(map[string]string),
		RemoveHeaders: []string{},
		AddHeaders:    make(map[string]string),
	}
}

// AddHeaderSet creates a middleware that sets a single header
func AddHeaderSet(key, value string) RequestHook {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// AddUserAgent creates a user agent middleware
func AddUserAgent(userAgent string) RequestHook {
	return AddHeaderSet("User-Agent", userAgent)
}

// AddContentType creates a content type middleware
func AddContentType(contentType string) RequestHook {
	return AddHeaderSet("Content-Type", contentType)
}

// AddAccept creates an accept header middleware
func AddAccept(accept string) RequestHook {
	return AddHeaderSet("Accept", accept)
}

// AddHeader creates a header hook based on configuration
func AddHeader(config HeaderConfig) RequestHook {
	return func(req *http.Request) error {
		if config.UserAgent != "" {
			req.Header.Set("User-Agent", config.UserAgent)
		}

		for key, value := range config.CustomHeaders {
			req.Header.Set(key, value)
		}

		for key, value := range config.AddHeaders {
			req.Header.Set(key, value)
		}

		for _, header := range config.RemoveHeaders {
			req.Header.Del(header)
		}

		return nil
	}
}

// AddSecurityHeaders creates a security headers middleware
func AddSecurityHeaders() RequestHook {
	return func(req *http.Request) error {
		req.Header.Set("X-Content-Type-Options", "nosniff")
		req.Header.Set("X-Frame-Options", "DENY")
		req.Header.Set("X-XSS-Protection", "1; mode=block")
		req.Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return nil
	}
}

// AddNoCache creates a no-cache middleware
func AddNoCache() RequestHook {
	return AddHeaderSet("Cache-Control", "no-cache, no-store, must-revalidate")
}
