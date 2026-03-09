package codewire

import (
	"context"
	"net/url"

	"codewire.sh/sdk-go/internal/api"
)

type SecretProjectService struct {
	client *Client
}

func (s *SecretProjectService) Create(ctx context.Context, req api.CreateSecretProjectInputBody) (*api.SecretProject, error) {
	path, err := s.client.orgPath("/secret-projects")
	if err != nil {
		return nil, err
	}
	var proj api.SecretProject
	if err := s.client.do(ctx, "POST", path, req, &proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

func (s *SecretProjectService) List(ctx context.Context) ([]api.SecretProject, error) {
	path, err := s.client.orgPath("/secret-projects")
	if err != nil {
		return nil, err
	}
	var projects []api.SecretProject
	if err := s.client.do(ctx, "GET", path, nil, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *SecretProjectService) Delete(ctx context.Context, projectID string) error {
	path, err := s.client.orgPath("/secret-projects/%s", projectID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *SecretProjectService) ListSecrets(ctx context.Context, projectID string) ([]api.SecretMetadata, error) {
	path, err := s.client.orgPath("/secret-projects/%s/secrets", projectID)
	if err != nil {
		return nil, err
	}
	var resp api.ListProjectSecretsOutputBody
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	if resp.Secrets == nil {
		return nil, nil
	}
	return *resp.Secrets, nil
}

func (s *SecretProjectService) SetSecret(ctx context.Context, projectID string, key, value string) error {
	path, err := s.client.orgPath("/secret-projects/%s/secrets", projectID)
	if err != nil {
		return err
	}
	req := api.SetProjectSecretInputBody{
		Key:   key,
		Value: value,
	}
	return s.client.do(ctx, "PUT", path, req, nil)
}

func (s *SecretProjectService) DeleteSecret(ctx context.Context, projectID, key string) error {
	path, err := s.client.orgPath("/secret-projects/%s/secrets/%s", projectID, url.PathEscape(key))
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}
