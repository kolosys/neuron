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

// UserAgentMiddleware creates a user agent middleware
func UserAgentMiddleware(userAgent string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// CustomHeaderMiddleware creates a custom header middleware
func CustomHeaderMiddleware(headers map[string]string) RequestMiddleware {
	return func(req *http.Request) error {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		return nil
	}
}

// HeaderMiddleware creates a header middleware based on configuration
func HeaderMiddleware(config HeaderConfig) RequestMiddleware {
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

// ContentTypeMiddleware creates a content type middleware
func ContentTypeMiddleware(contentType string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Content-Type", contentType)
		return nil
	}
}

// AcceptMiddleware creates an accept header middleware
func AcceptMiddleware(accept string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Accept", accept)
		return nil
	}
}

// JSONContentTypeMiddleware creates a JSON content type middleware
func JSONContentTypeMiddleware() RequestMiddleware {
	return ContentTypeMiddleware("application/json")
}

// XMLContentTypeMiddleware creates an XML content type middleware
func XMLContentTypeMiddleware() RequestMiddleware {
	return ContentTypeMiddleware("application/xml")
}

// FormContentTypeMiddleware creates a form content type middleware
func FormContentTypeMiddleware() RequestMiddleware {
	return ContentTypeMiddleware("application/x-www-form-urlencoded")
}

// MultipartContentTypeMiddleware creates a multipart content type middleware
func MultipartContentTypeMiddleware() RequestMiddleware {
	return ContentTypeMiddleware("multipart/form-data")
}

// AcceptJSONMiddleware creates an accept JSON middleware
func AcceptJSONMiddleware() RequestMiddleware {
	return AcceptMiddleware("application/json")
}

// AcceptXMLMiddleware creates an accept XML middleware
func AcceptXMLMiddleware() RequestMiddleware {
	return AcceptMiddleware("application/xml")
}

// AcceptAllMiddleware creates an accept all middleware
func AcceptAllMiddleware() RequestMiddleware {
	return AcceptMiddleware("*/*")
}

// CORSHeadersMiddleware creates a CORS headers middleware
func CORSHeadersMiddleware() RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Origin", req.Header.Get("Origin"))
		req.Header.Set("Access-Control-Request-Method", req.Method)
		req.Header.Set("Access-Control-Request-Headers", strings.Join(getHeaderNames(req.Header), ", "))
		return nil
	}
}

// SecurityHeadersMiddleware creates a security headers middleware
func SecurityHeadersMiddleware() RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("X-Content-Type-Options", "nosniff")
		req.Header.Set("X-Frame-Options", "DENY")
		req.Header.Set("X-XSS-Protection", "1; mode=block")
		req.Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return nil
	}
}

// CacheControlMiddleware creates a cache control middleware
func CacheControlMiddleware(cacheControl string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("Cache-Control", cacheControl)
		return nil
	}
}

// NoCacheMiddleware creates a no-cache middleware
func NoCacheMiddleware() RequestMiddleware {
	return CacheControlMiddleware("no-cache, no-store, must-revalidate")
}

// ETagMiddleware creates an ETag middleware
func ETagMiddleware(etag string) RequestMiddleware {
	return func(req *http.Request) error {
		req.Header.Set("If-None-Match", etag)
		return nil
	}
}

// ConditionalRequestMiddleware creates a conditional request middleware
func ConditionalRequestMiddleware(ifModifiedSince, ifNoneMatch string) RequestMiddleware {
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

// HeaderTransformationMiddleware creates a header transformation middleware
func HeaderTransformationMiddleware(transformFunc func(string, string) (string, string)) RequestMiddleware {
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

// HeaderFilterMiddleware creates a header filter middleware
func HeaderFilterMiddleware(allowedHeaders []string) RequestMiddleware {
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
