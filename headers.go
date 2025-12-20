package neuron

import (
	"net/http"
	"strings"
)

// HeaderConfig configures header middleware
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

// AddUserAgent creates a user agent middleware
func AddUserAgent(userAgent string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// AddCustomHeader creates a custom header middleware
func AddCustomHeader(headers map[string]string) RequestMiddleware {
	return func(req *http.Request) error {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		return nil
	}
}

// AddHeader creates a header middleware based on configuration
func AddHeader(config HeaderConfig) RequestMiddleware {
	return func(req *http.Request) error {
		// Set user agent
		if config.UserAgent != "" {
			req.Header.Set("User-Agent", config.UserAgent)
		}

		// Add custom headers
		for key, value := range config.CustomHeaders {
			req.Header.Set(key, value)
		}

		// Add additional headers
		for key, value := range config.AddHeaders {
			req.Header.Set(key, value)
		}

		// Remove specified headers
		for _, header := range config.RemoveHeaders {
			req.Header.Del(header)
		}

		return nil
	}
}

// AddContentType creates a content type middleware
func AddContentType(contentType string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Content-Type", contentType)
		return nil
	}
}

// AddAccept creates an accept header middleware
func AddAccept(accept string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Accept", accept)
		return nil
	}
}

// AddJSONContentType creates a JSON content type middleware
func AddJSONContentType() RequestMiddleware {
	return AddContentType("application/json")
}

// AddXMLContentType creates an XML content type middleware
func AddXMLContentType() RequestMiddleware {
	return AddContentType("application/xml")
}

// AddFormContentType creates a form content type middleware
func AddFormContentType() RequestMiddleware {
	return AddContentType("application/x-www-form-urlencoded")
}

// AddMultipartContentType creates a multipart content type middleware
func AddMultipartContentType() RequestMiddleware {
	return AddContentType("multipart/form-data")
}

// AddAcceptJSON creates an accept JSON middleware
func AddAcceptJSON() RequestMiddleware {
	return AddAccept("application/json")
}

// AddAcceptXML creates an accept XML middleware
func AddAcceptXML() RequestMiddleware {
	return AddAccept("application/xml")
}

// AddAcceptAll creates an accept all middleware
func AddAcceptAll() RequestMiddleware {
	return AddAccept("*/*")
}

// AddCORSHeaders creates a CORS headers middleware
func AddCORSHeaders() RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Origin", req.Header.Get("Origin"))
		req.Header.Set("Access-Control-Request-Method", req.Method)
		req.Header.Set("Access-Control-Request-Headers", strings.Join(getHeaderNames(req.Header), ", "))
		return nil
	}
}

// AddSecurityHeaders creates a security headers middleware
func AddSecurityHeaders() RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("X-Content-Type-Options", "nosniff")
		req.Header.Set("X-Frame-Options", "DENY")
		req.Header.Set("X-XSS-Protection", "1; mode=block")
		req.Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return nil
	}
}

// AddCacheControl creates a cache control middleware
func AddCacheControl(cacheControl string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Cache-Control", cacheControl)
		return nil
	}
}

// AddNoCache creates a no-cache middleware
func AddNoCache() RequestMiddleware {
	return AddCacheControl("no-cache, no-store, must-revalidate")
}

// AddETag creates an ETag middleware
func AddETag(etag string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("If-None-Match", etag)
		return nil
	}
}

// AddConditionalRequest creates a conditional request middleware
func AddConditionalRequest(ifModifiedSince, ifNoneMatch string) RequestMiddleware {
	return func(req *http.Request) error {
		if ifModifiedSince != "" {
			req.Header.Set("If-Modified-Since", ifModifiedSince)
		}
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}
		return nil
	}
}

// AddHeaderTransformation creates a header transformation middleware
func AddHeaderTransformation(transformFunc func(string, string) (string, string)) RequestMiddleware {
	return func(req *http.Request) error {
		newHeaders := make(map[string]string)

		for key, values := range req.Header {
			for _, value := range values {
				newKey, newValue := transformFunc(key, value)
				newHeaders[newKey] = newValue
			}
		}

		// Clear existing headers
		for key := range req.Header {
			req.Header.Del(key)
		}

		// Set transformed headers
		for key, value := range newHeaders {
			req.Header.Set(key, value)
		}

		return nil
	}
}

// AddHeaderFilter creates a header filter middleware
func AddHeaderFilter(allowedHeaders []string) RequestMiddleware {
	return func(req *http.Request) error {
		// Create a map of allowed headers for quick lookup
		allowed := make(map[string]bool)
		for _, header := range allowedHeaders {
			allowed[strings.ToLower(header)] = true
		}

		// Remove headers that are not allowed
		for key := range req.Header {
			if !allowed[strings.ToLower(key)] {
				req.Header.Del(key)
			}
		}

		return nil
	}
}

// getHeaderNames returns a slice of header names
func getHeaderNames(headers http.Header) []string {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	return names
}
