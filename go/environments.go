package codewire

import (
	"context"
	"fmt"
	"net/url"
)

type EnvironmentService struct {
	client *Client
}

func (s *EnvironmentService) Create(ctx context.Context, orgID string, req CreateEnvironmentRequest) (*Environment, error) {
	var env Environment
	err := s.client.do(ctx, "POST", fmt.Sprintf("/organizations/%s/environments", orgID), req, &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (s *EnvironmentService) List(ctx context.Context, orgID string, params *ListEnvironmentsParams) ([]Environment, error) {
	path := fmt.Sprintf("/organizations/%s/environments", orgID)
	if params != nil {
		q := url.Values{}
		if params.Type != "" {
			q.Set("type", params.Type)
		}
		if params.State != "" {
			q.Set("state", params.State)
		}
		if qs := q.Encode(); qs != "" {
			path += "?" + qs
		}
	}
	var envs []Environment
	err := s.client.do(ctx, "GET", path, nil, &envs)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

func (s *EnvironmentService) Get(ctx context.Context, orgID, envID string) (*Environment, error) {
	var env Environment
	err := s.client.do(ctx, "GET", fmt.Sprintf("/organizations/%s/environments/%s", orgID, envID), nil, &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (s *EnvironmentService) Delete(ctx context.Context, orgID, envID string) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/organizations/%s/environments/%s", orgID, envID), nil, nil)
}

func (s *EnvironmentService) Stop(ctx context.Context, orgID, envID string) (*StatusResponse, error) {
	var resp StatusResponse
	err := s.client.do(ctx, "POST", fmt.Sprintf("/organizations/%s/environments/%s/stop", orgID, envID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *EnvironmentService) Start(ctx context.Context, orgID, envID string) (*StatusResponse, error) {
	var resp StatusResponse
	err := s.client.do(ctx, "POST", fmt.Sprintf("/organizations/%s/environments/%s/start", orgID, envID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
