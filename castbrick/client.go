package castbrick

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.castbrick.co/v1"

	// defaultTimeout bounds every request made with the SDK's default HTTP
	// client, so a slow or unresponsive upstream cannot pin goroutines forever.
	defaultTimeout = 30 * time.Second

	// maxResponseBytes caps how much of a response body is read into memory.
	maxResponseBytes = 10 << 20 // 10 MiB

	// maxErrorBodyBytes caps how much of an error response is kept in APIError.Body.
	maxErrorBodyBytes = 64 << 10 // 64 KiB

	// maxErrorMessageBytes caps how much of the body is rendered by APIError.Error(),
	// keeping log lines small and limiting leakage of response data into logs.
	maxErrorMessageBytes = 512

	// maxIDLength is the longest resource identifier accepted in a URL path.
	maxIDLength = 256

	userAgent = "castbrick-go"
)

var (
	// ErrInvalidConfig is returned by every call made with a client whose API
	// key or base URL is invalid (for example an empty key or a non-HTTPS URL).
	ErrInvalidConfig = errors.New("castbrick: invalid client configuration")

	// ErrInvalidID is returned when a resource identifier is empty, too long or
	// contains characters that could alter the request path or query.
	ErrInvalidID = errors.New("castbrick: invalid resource ID")

	// ErrResponseTooLarge is returned when the API response exceeds the size limit.
	ErrResponseTooLarge = errors.New("castbrick: response body too large")
)

// Client is the low-level HTTP client used by all resources.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
	// configErr is set when the client was built with invalid settings. Every
	// request fails with it, so a misconfiguration can never silently send the
	// API key somewhere unsafe.
	configErr error
}

func newClient(apiKey, baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout:       defaultTimeout,
			CheckRedirect: checkRedirect,
		}
	}
	c := &Client{apiKey: apiKey, http: httpClient}

	normalized, err := validateBaseURL(baseURL)
	if err != nil {
		c.configErr = err
		return c
	}
	c.baseURL = normalized

	if strings.TrimSpace(apiKey) == "" {
		c.configErr = fmt.Errorf("%w: API key is empty", ErrInvalidConfig)
	}
	return c
}

// validateBaseURL requires HTTPS, except for loopback hosts used in local
// development and tests, and rejects URLs carrying credentials, a query or a fragment.
func validateBaseURL(raw string) (string, error) {
	if raw == "" {
		raw = defaultBaseURL
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("%w: base URL: %v", ErrInvalidConfig, err)
	}
	if u.Host == "" {
		return "", fmt.Errorf("%w: base URL %q has no host", ErrInvalidConfig, raw)
	}
	if u.User != nil {
		return "", fmt.Errorf("%w: base URL must not contain credentials", ErrInvalidConfig)
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("%w: base URL must not contain a query or fragment", ErrInvalidConfig)
	}
	switch u.Scheme {
	case "https":
	case "http":
		if !isLoopback(u.Hostname()) {
			return "", fmt.Errorf("%w: base URL must use https (http is only allowed for localhost)", ErrInvalidConfig)
		}
	default:
		return "", fmt.Errorf("%w: unsupported base URL scheme %q", ErrInvalidConfig, u.Scheme)
	}
	return strings.TrimRight(u.String(), "/"), nil
}

func isLoopback(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// checkRedirect stops the default client from following a redirect that
// changes host or scheme, so the Authorization header is never forwarded to
// another origin or downgraded to plain HTTP.
func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("castbrick: stopped after 10 redirects")
	}
	orig := via[0].URL
	if req.URL.Scheme != orig.Scheme || !strings.EqualFold(req.URL.Host, orig.Host) {
		return fmt.Errorf("castbrick: refusing redirect from %s://%s to %s://%s",
			orig.Scheme, orig.Host, req.URL.Scheme, req.URL.Host)
	}
	return nil
}

// escapeID validates a resource identifier and escapes it for use as a single
// URL path segment. It rejects values that could traverse to another endpoint
// (".." or "/") or inject a query or fragment ("?", "#"), even when escaped,
// because some proxies and servers decode %2F before routing.
func escapeID(id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidID)
	}
	if len(id) > maxIDLength {
		return "", fmt.Errorf("%w: longer than %d characters", ErrInvalidID, maxIDLength)
	}
	if id == "." || id == ".." {
		return "", fmt.Errorf("%w: %q", ErrInvalidID, id)
	}
	for _, r := range id {
		if r < 0x20 || r == 0x7f || strings.ContainsRune(`/\?#%`, r) {
			return "", fmt.Errorf("%w: contains a forbidden character", ErrInvalidID)
		}
	}
	return url.PathEscape(id), nil
}

// pathf builds a request path from a format string whose %s verbs are filled,
// in order, with validated and escaped resource identifiers.
func pathf(format string, ids ...string) (string, error) {
	args := make([]any, len(ids))
	for i, id := range ids {
		esc, err := escapeID(id)
		if err != nil {
			return "", err
		}
		args[i] = esc
	}
	return fmt.Sprintf(format, args...), nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	if c.configErr != nil {
		return nil, c.configErr
	}

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return err
	}
	if len(respBody) > maxResponseBytes {
		return fmt.Errorf("%w: exceeds %d bytes", ErrResponseTooLarge, maxResponseBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, respBody)
	}

	if out != nil && resp.StatusCode != http.StatusNoContent && len(respBody) > 0 {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	return c.getQuery(ctx, path, nil, out)
}

func (c *Client) getQuery(ctx context.Context, path string, query url.Values, out any) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) put(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) delete(ctx context.Context, path string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	// Message is the human-readable error extracted from the JSON body, if any.
	Message string
	// Body is the raw response body, truncated to 64 KiB.
	Body string
}

func newAPIError(status int, body []byte) *APIError {
	if len(body) > maxErrorBodyBytes {
		body = body[:maxErrorBodyBytes]
	}
	return &APIError{StatusCode: status, Message: extractMessage(body), Body: string(body)}
}

// extractMessage reads the common error fields used by the API, including
// ASP.NET ProblemDetails ("detail", "title").
func extractMessage(body []byte) string {
	var payload map[string]any
	if json.Unmarshal(body, &payload) != nil {
		return ""
	}
	for _, key := range []string{"message", "error", "detail", "title"} {
		if s, ok := payload[key].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.Body
	}
	return fmt.Sprintf("castbrick: API error %d: %s", e.StatusCode, sanitizeForLog(msg, maxErrorMessageBytes))
}

// sanitizeForLog truncates s and replaces control characters so an error
// returned by the API cannot forge extra log lines.
func sanitizeForLog(s string, max int) string {
	truncated := false
	if len(s) > max {
		s = s[:max]
		truncated = true
	}
	s = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, strings.ToValidUTF8(s, "?"))
	if truncated {
		s += "…(truncated)"
	}
	return s
}
