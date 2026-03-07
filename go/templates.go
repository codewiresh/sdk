package codewire

import (
	"context"
	"fmt"
	"net/url"
)

type TemplateService struct {
	client *Client
}

func (s *TemplateService) Create(ctx context.Context, orgID string, req CreateTemplateRequest) (*Template, error) {
	var tmpl Template
	err := s.client.do(ctx, "POST", fmt.Sprintf("/organizations/%s/templates", orgID), req, &tmpl)
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (s *TemplateService) List(ctx context.Context, orgID string, params *ListTemplatesParams) ([]Template, error) {
	path := fmt.Sprintf("/organizations/%s/templates", orgID)
	if params != nil {
		q := url.Values{}
		if params.Type != "" {
			q.Set("type", params.Type)
		}
		if qs := q.Encode(); qs != "" {
			path += "?" + qs
		}
	}
	var templates []Template
	err := s.client.do(ctx, "GET", path, nil, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (s *TemplateService) Delete(ctx context.Context, orgID, templateID string) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/organizations/%s/templates/%s", orgID, templateID), nil, nil)
}
