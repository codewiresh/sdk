package codewire

import (
	"context"

	"codewire.sh/sdk-go/internal/api"
)

type APIKeyService struct {
	client *Client
}

func (s *APIKeyService) Create(ctx context.Context, req api.CreateAPIKeyBody) (*api.CreateAPIKeyResponse, error) {
	path, err := s.client.orgPath("/api-keys")
	if err != nil {
		return nil, err
	}
	var resp api.CreateAPIKeyResponse
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *APIKeyService) List(ctx context.Context) ([]api.APIKey, error) {
	path, err := s.client.orgPath("/api-keys")
	if err != nil {
		return nil, err
	}
	var keys []api.APIKey
	if err := s.client.do(ctx, "GET", path, nil, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *APIKeyService) Delete(ctx context.Context, keyID string) error {
	path, err := s.client.orgPath("/api-keys/%s", keyID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}
