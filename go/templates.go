package codewire

import (
	"context"
	"net/url"

	"codewire.sh/sdk-go/internal/api"
)

type TemplateService struct {
	client *Client
}

func (s *TemplateService) Create(ctx context.Context, req api.CreateTemplateBody) (*api.EnvironmentTemplate, error) {
	path, err := s.client.orgPath("/templates")
	if err != nil {
		return nil, err
	}
	var tmpl api.EnvironmentTemplate
	if err := s.client.do(ctx, "POST", path, req, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (s *TemplateService) List(ctx context.Context, params *api.ListEnvironmentTemplatesParams) ([]api.EnvironmentTemplate, error) {
	path, err := s.client.orgPath("/templates")
	if err != nil {
		return nil, err
	}
	if params != nil && params.Type != nil {
		q := url.Values{}
		q.Set("type", *params.Type)
		path += "?" + q.Encode()
	}
	var templates []api.EnvironmentTemplate
	if err := s.client.do(ctx, "GET", path, nil, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (s *TemplateService) Get(ctx context.Context, templateID string) (*api.EnvironmentTemplate, error) {
	path, err := s.client.orgPath("/templates/%s", templateID)
	if err != nil {
		return nil, err
	}
	var tmpl api.EnvironmentTemplate
	if err := s.client.do(ctx, "GET", path, nil, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (s *TemplateService) Update(ctx context.Context, templateID string, req api.UpdateTemplateBody) (*api.EnvironmentTemplate, error) {
	path, err := s.client.orgPath("/templates/%s", templateID)
	if err != nil {
		return nil, err
	}
	var tmpl api.EnvironmentTemplate
	if err := s.client.do(ctx, "PATCH", path, req, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (s *TemplateService) Delete(ctx context.Context, templateID string) error {
	path, err := s.client.orgPath("/templates/%s", templateID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}
