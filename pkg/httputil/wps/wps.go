package wps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultContentType  = "application/json"
	defaultTimeout      = 30 * time.Second
	AuthorizationHeader = "Authorization"
	RFC1123             = "Mon, 02 Jan 2006 15:04:05 GMT"
	// RFC1123             = "Monday, 02 Jan 2006 15:04:05 MST"
)

var (
	ErrInvalidRequest = errors.New("invalid request parameters")
	ErrHTTPRequest    = errors.New("HTTP request failed")
	openApiPathPrefix = "/openapi"
)

type WPSRequest struct {
	baseURL     string
	method      string
	path        string
	body        []byte
	contentType string
	ksoDate     string
	timeout     time.Duration
	headers     map[string]string
	queryParams map[string]string
	accessKey   string
	secretKey   string
	client      *http.Client
	accessToken string
}

type Option func(*WPSRequest)

func NewWPSRequest(baseURL, accessKey, secretKey string, opts ...Option) *WPSRequest {
	r := &WPSRequest{
		baseURL:     strings.TrimRight(baseURL, "/"),
		accessKey:   accessKey,
		secretKey:   secretKey,
		timeout:     defaultTimeout,
		contentType: defaultContentType,
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
		client:      &http.Client{Timeout: defaultTimeout},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Option setters
func WithMethod(method string) Option {
	return func(r *WPSRequest) {
		r.method = strings.ToUpper(method)
	}
}

func WithPath(path string) Option {
	return func(r *WPSRequest) {
		//r.path = strings.TrimLeft(path, "/")
		r.path = path
	}
}

func WithBody(body []byte) Option {
	return func(r *WPSRequest) {
		r.body = body
	}
}

func WithJSONBody(v interface{}) Option {
	return func(r *WPSRequest) {
		if b, err := json.Marshal(v); err == nil {
			r.body = b
			r.contentType = "application/json"
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(r *WPSRequest) {
		r.timeout = timeout
		r.client.Timeout = timeout
	}
}

func WithHeader(key, value string) Option {
	return func(r *WPSRequest) {
		r.headers[key] = value
	}
}

func WithQueryParam(key, value string) Option {
	return func(r *WPSRequest) {
		r.queryParams[key] = value
	}
}

func WithContentType(contentType string) Option {
	return func(r *WPSRequest) {
		r.contentType = contentType
	}
}

func WithKsoDate(date string) Option {
	return func(r *WPSRequest) {
		r.ksoDate = date
	}
}

func WithAuthorization(accessToken string) Option {
	return func(r *WPSRequest) {
		r.accessToken = "Bearer " + accessToken
	}
}

// Core methods
func (r *WPSRequest) BuildRequest() (*http.Request, error) {
	if r.method == "" || r.path == "" {
		return nil, fmt.Errorf("%w: method and path are required", ErrInvalidRequest)
	}

	// Build URL with query params
	u, err := url.Parse(fmt.Sprintf("%s%s", r.baseURL, r.path))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid URL", ErrInvalidRequest)
	}

	q := u.Query()
	for k, v := range r.queryParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest(r.method, u.String(), bytes.NewBuffer(r.body))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request", ErrHTTPRequest)
	}

	// Set headers
	if r.body != nil && r.contentType != "" {
		req.Header.Set("Content-Type", r.contentType)
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (r *WPSRequest) Do(ctx context.Context) ([]byte, error) {
	req, err := r.BuildRequest()
	if err != nil {
		return nil, err
	}

	// Add KSO signature
	signer, err := NewKsoSign(r.accessKey, r.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	signPath := strings.TrimPrefix(req.URL.Path, openApiPathPrefix)
	sign, err := signer.KSO1Sign(r.method, signPath, r.contentType, r.ksoDate, r.body)

	// fmt.Printf("out: %v\n", sign)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	req.Header.Set(KsoDateHeader, sign.Date)
	req.Header.Set(KsoAuthHeader, sign.Authorization)
	req.Header.Set(AuthorizationHeader, r.accessToken)

	// for k, v := range req.Header {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
	// Execute request
	// Execute request
	resp, err := r.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("%w: status %d", ErrHTTPRequest, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Convenience method
func (r *WPSRequest) PostJSON(ctx context.Context, path string, accessToken string, body interface{}) ([]byte, error) {
	req := NewWPSRequest(r.baseURL, r.accessKey, r.secretKey,
		WithMethod(http.MethodPost),
		WithPath(path),
		WithJSONBody(body),
		WithKsoDate(time.Now().UTC().Format(RFC1123)),
		WithAuthorization(accessToken),
	)
	return req.Do(ctx)
}
