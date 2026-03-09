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
	Environments   *EnvironmentService
	Templates      *TemplateService
	APIKeys        *APIKeyService
	Secrets        *SecretService
	SecretProjects *SecretProjectService

	baseURL    string
	apiKey     string
	orgID      string
	httpClient *http.Client
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(url, "/") }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

func WithOrgID(orgID string) Option {
	return func(c *Client) { c.orgID = orgID }
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
	if c.orgID == "" {
		c.orgID = os.Getenv("CODEWIRE_ORG_ID")
	}
	c.Environments = &EnvironmentService{client: c}
	c.Templates = &TemplateService{client: c}
	c.APIKeys = &APIKeyService{client: c}
	c.Secrets = &SecretService{client: c}
	c.SecretProjects = &SecretProjectService{client: c}
	return c
}

func (c *Client) resolveOrgID() (string, error) {
	if c.orgID != "" {
		return c.orgID, nil
	}
	return "", fmt.Errorf("codewire: orgID is required — use WithOrgID or set CODEWIRE_ORG_ID")
}

func (c *Client) orgPath(pattern string, args ...any) (string, error) {
	orgID, err := c.resolveOrgID()
	if err != nil {
		return "", err
	}
	allArgs := make([]any, 0, len(args)+1)
	allArgs = append(allArgs, orgID)
	allArgs = append(allArgs, args...)
	return fmt.Sprintf("/organizations/%s"+pattern, allArgs...), nil
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

func (c *Client) doRaw(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("codewire: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("codewire: request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
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
		return nil, &Error{StatusCode: resp.StatusCode, Detail: detail}
	}

	return resp, nil
}
