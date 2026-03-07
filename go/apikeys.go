package codewire

import (
	"context"
	"fmt"
)

type APIKeyService struct {
	client *Client
}

func (s *APIKeyService) Create(ctx context.Context, orgID string, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	var resp CreateAPIKeyResponse
	err := s.client.do(ctx, "POST", fmt.Sprintf("/organizations/%s/api-keys", orgID), req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *APIKeyService) List(ctx context.Context, orgID string) ([]APIKey, error) {
	var keys []APIKey
	err := s.client.do(ctx, "GET", fmt.Sprintf("/organizations/%s/api-keys", orgID), nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *APIKeyService) Delete(ctx context.Context, orgID, keyID string) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/organizations/%s/api-keys/%s", orgID, keyID), nil, nil)
}
