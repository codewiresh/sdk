package codewire

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DefaultBaseURL = "https://api.codewire.sh"

type Client struct {
	Environments *EnvironmentService
	Templates    *TemplateService
	APIKeys      *APIKeyService

	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(url, "/") }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

func New(apiKey string, opts ...Option) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("CODEWIRE_API_KEY")
	}
	c := &Client{
		baseURL:    DefaultBaseURL,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	c.Environments = &EnvironmentService{client: c}
	c.Templates = &TemplateService{client: c}
	c.APIKeys = &APIKeyService{client: c}
	return c
}

type Error struct {
	StatusCode int
	Detail     string
}

func (e *Error) Error() string {
	return fmt.Sprintf("codewire: %d %s", e.StatusCode, e.Detail)
}

func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("codewire: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("codewire: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("codewire: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		detail := resp.Status
		var errBody struct {
			Detail  string `json:"detail"`
			Message string `json:"message"`
		}
		if json.NewDecoder(resp.Body).Decode(&errBody) == nil {
			if errBody.Detail != "" {
				detail = errBody.Detail
			} else if errBody.Message != "" {
				detail = errBody.Message
			}
		}
		return &Error{StatusCode: resp.StatusCode, Detail: detail}
	}

	if resp.StatusCode == 204 || result == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
