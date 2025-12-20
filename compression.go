package neuron

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// CompressionType represents the type of compression
type CompressionType int

const (
	CompressionTypeGzip CompressionType = iota
	CompressionTypeDeflate
	CompressionTypeBrotli
)

// CompressionConfig configures compression middleware
type CompressionConfig struct {
	Type         CompressionType
	Level        int
	MinSize      int
	ContentTypes []string
}

// DefaultCompressionConfig returns a default compression configuration
func DefaultCompressionConfig() CompressionConfig {
	return CompressionConfig{
		Type:         CompressionTypeGzip,
		Level:        6,    // Default compression level
		MinSize:      1024, // 1KB minimum size
		ContentTypes: []string{"application/json", "text/plain", "text/html", "application/xml"},
	}
}

// AddCompression creates a request compression middleware
func AddCompression(config CompressionConfig) RequestHook {
	return func(req *http.Request) error {
		// Check if request body should be compressed
		if req.Body == nil {
			return nil
		}

		// Check content type
		contentType := req.Header.Get("Content-Type")
		if !shouldCompress(contentType, config.ContentTypes) {
			return nil
		}

		// Add compression headers
		switch config.Type {
		case CompressionTypeGzip:
			req.Header.Set("Content-Encoding", "gzip")
		case CompressionTypeDeflate:
			req.Header.Set("Content-Encoding", "deflate")
		case CompressionTypeBrotli:
			req.Header.Set("Content-Encoding", "br")
		}

		// Set accept encoding for response
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")

		return nil
	}
}

// AddResponseDecompression creates a response decompression middleware
func AddResponseDecompression() ResponseHook {
	return func(resp *http.Response) error {
		// Check if response is compressed
		encoding := resp.Header.Get("Content-Encoding")
		if encoding == "" {
			return nil
		}

		// Decompress based on encoding
		switch encoding {
		case "gzip":
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			resp.Body = &gzipReader{reader: reader, original: resp.Body}
		case "deflate":
			// Note: deflate decompression would need additional implementation
			// This is a placeholder
		case "br":
			// Note: brotli decompression would need additional implementation
			// This is a placeholder
		}

		return nil
	}
}

// gzipReader wraps a gzip reader to properly close the underlying reader
type gzipReader struct {
	reader   *gzip.Reader
	original io.ReadCloser
}

func (r *gzipReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *gzipReader) Close() error {
	if err := r.reader.Close(); err != nil {
		return err
	}
	return r.original.Close()
}

// shouldCompress checks if content should be compressed
func shouldCompress(contentType string, allowedTypes []string) bool {
	if contentType == "" {
		return false
	}

	for _, allowedType := range allowedTypes {
		if strings.Contains(contentType, allowedType) {
			return true
		}
	}

	return false
}

// AddResponseCompression creates a response compression middleware
func AddResponseCompression(config CompressionConfig) ResponseHook {
	return func(resp *http.Response) error {
		// Check if response should be compressed
		contentType := resp.Header.Get("Content-Type")
		if !shouldCompress(contentType, config.ContentTypes) {
			return nil
		}

		// Check if client accepts compression
		acceptEncoding := resp.Request.Header.Get("Accept-Encoding")
		if acceptEncoding == "" {
			return nil
		}

		// Add compression headers
		switch config.Type {
		case CompressionTypeGzip:
			if strings.Contains(acceptEncoding, "gzip") {
				resp.Header.Set("Content-Encoding", "gzip")
			}
		case CompressionTypeDeflate:
			if strings.Contains(acceptEncoding, "deflate") {
				resp.Header.Set("Content-Encoding", "deflate")
			}
		case CompressionTypeBrotli:
			if strings.Contains(acceptEncoding, "br") {
				resp.Header.Set("Content-Encoding", "br")
			}
		}

		return nil
	}
}

// AddAutoCompression creates an automatic compression middleware
func AddAutoCompression() RequestHook {
	return func(req *http.Request) error {
		// Set accept encoding for all requests
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		return nil
	}
}

// AddCompressionLevel creates a compression level middleware
func AddCompressionLevel(level int) RequestHook {
	return func(req *http.Request) error {
		// This would be used with a custom HTTP client that supports compression levels
		// For now, it's a placeholder for future implementation
		return nil
	}
}
