package neuron

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Stream sends an HTTP request and returns the raw response with Body still
// open. The caller must Close the body.
//
// Unlike Execute, Stream never retries after headers or a body exist —
// including 5xx. Closing and re-POSTing an LLM stream can duplicate tool
// calls. Transport errors that occur before any response are retried
// according to MaxRetries / RequestOptions.Retries.
//
// Stream does not apply AdaptiveTimeout. LLM callers should set
// http.Client.Timeout (around 10 minutes) or a request context deadline.
func Stream[TRequest any](client *Client, method HTTPMethod, path string, request TRequest, options *RequestOptions) (*http.Response, error) {
	if options == nil {
		options = &RequestOptions{}
	}
	if options.Context == nil {
		options.Context = context.Background()
	}

	req, err := client.buildRequest(method, path, request, options)
	if err != nil {
		ce := ClientError{
			Type:      ErrorTypeRequest,
			Message:   fmt.Sprintf("failed to build request: %v", err),
			Route:     path,
			Cause:     err,
			Timestamp: time.Now(),
		}
		if req != nil {
			ce = ce.WithContext(req, 0)
		}
		return nil, ce
	}

	startTime := time.Now()
	req = req.WithContext(context.WithValue(req.Context(), requestStartKey, startTime))

	requestHooks := append(client.Config.RequestHooks, options.RequestHooks...)
	responseHooks := append(client.Config.ResponseHooks, options.ResponseHooks...)

	for _, hook := range requestHooks {
		if err := hook(req); err != nil {
			if ce, ok := err.(ClientError); ok {
				return nil, ce
			}
			return nil, ClientError{
				Type:    ErrorTypeRequest,
				Message: fmt.Sprintf("request hook failed: %v", err),
				Route:   path,
				Cause:   err,
			}.WithContext(req, 0)
		}
	}

	resp, err := client.streamRequest(req, options)
	if err != nil {
		return nil, err
	}

	client.Metrics.RecordResponse(resp.StatusCode, time.Since(startTime))

	for _, hook := range responseHooks {
		if err := hook(resp); err != nil {
			resp.Body.Close()
			if ce, ok := err.(ClientError); ok {
				return nil, ce
			}
			return nil, ClientError{
				Type:    ErrorTypeResponse,
				Message: fmt.Sprintf("response hook failed: %v", err),
				Route:   path,
				Cause:   err,
			}.WithContext(req, 0)
		}
	}

	return resp, nil
}

func (c *Client) streamRequest(req *http.Request, options *RequestOptions) (*http.Response, error) {
	if c.Config.RateLimiter != nil {
		if err := c.Config.RateLimiter.Wait(req.Context(), req.Method, req.URL.Path); err != nil {
			return nil, ClientError{
				Type:    ErrorTypeRateLimit,
				Message: "rate limit wait failed",
				Route:   req.URL.Path,
				Cause:   err,
			}.WithContext(req, 0)
		}
	}
	return c.streamWithRetries(req, options)
}

func (c *Client) streamWithRetries(req *http.Request, options *RequestOptions) (*http.Response, error) {
	maxRetries := c.Config.MaxRetries
	if options.Retries != nil {
		maxRetries = *options.Retries
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-req.Context().Done():
			return nil, ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled or timed out",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
		default:
		}

		resp, err := c.doHTTP(req)
		if err == nil {
			c.updateRateLimitFromResponse(req, resp)
			c.Metrics.RequestCount.Add(1)
			return resp, nil
		}

		lastErr = err
		if ce, ok := circuitClientError(err); ok {
			c.Metrics.RequestCount.Add(1)
			c.Metrics.ErrorCount.Add(1)
			return nil, ce.WithContext(req, attempt)
		}

		if attempt >= maxRetries {
			break
		}

		delay := c.calculateBackoffDelayWithJitter(attempt)
		select {
		case <-time.After(delay):
		case <-req.Context().Done():
			return nil, ClientError{
				Type:    ErrorTypeTimeout,
				Message: "request cancelled during retry backoff",
				Cause:   req.Context().Err(),
			}.WithContext(req, attempt)
		}
	}

	c.Metrics.RequestCount.Add(1)
	c.Metrics.ErrorCount.Add(1)
	return nil, ClientError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("request failed after %d attempts: %v", maxRetries+1, lastErr),
		Cause:   lastErr,
	}.WithContext(req, maxRetries)
}
