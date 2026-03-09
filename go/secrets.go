package codewire

import (
	"context"
	"fmt"
	"net/url"

	"codewire.sh/sdk-go/internal/api"
)

type SecretService struct {
	client *Client
}

func (s *SecretService) List(ctx context.Context) ([]api.SecretMetadata, error) {
	orgID, err := s.client.resolveOrgID()
	if err != nil {
		return nil, err
	}
	path := "/secrets?org_id=" + url.QueryEscape(orgID)
	var resp api.ListSecretsOutputBody
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Secrets == nil {
		return nil, nil
	}
	return *resp.Secrets, nil
}

func (s *SecretService) Set(ctx context.Context, key, value string) error {
	orgID, err := s.client.resolveOrgID()
	if err != nil {
		return err
	}
	req := api.SetSecretInputBody{
		OrgId: orgID,
		Key:   key,
		Value: value,
	}
	return s.client.do(ctx, "PUT", "/secrets", req, nil)
}

func (s *SecretService) Delete(ctx context.Context, key string) error {
	orgID, err := s.client.resolveOrgID()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/secrets/%s?org_id=%s", url.PathEscape(key), url.QueryEscape(orgID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *SecretService) ListUser(ctx context.Context) ([]api.SecretMetadata, error) {
	var resp api.ListSecretsOutputBody
	if err := s.client.do(ctx, "GET", "/user/secrets", nil, &resp); err != nil {
		return nil, err
	}
	if resp.Secrets == nil {
		return nil, nil
	}
	return *resp.Secrets, nil
}

func (s *SecretService) SetUser(ctx context.Context, key, value string) error {
	req := api.SetUserSecretInputBody{
		Key:   key,
		Value: value,
	}
	return s.client.do(ctx, "PUT", "/user/secrets", req, nil)
}

func (s *SecretService) DeleteUser(ctx context.Context, key string) error {
	path := fmt.Sprintf("/user/secrets/%s", url.PathEscape(key))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}
