package codewire

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"codewire.sh/sdk-go/internal/api"
)

type EnvironmentService struct {
	client *Client
}

func (s *EnvironmentService) Create(ctx context.Context, req api.CreateEnvironmentBody) (*Environment, error) {
	path, err := s.client.orgPath("/environments")
	if err != nil {
		return nil, err
	}
	var data api.Environment
	if err := s.client.do(ctx, "POST", path, req, &data); err != nil {
		return nil, err
	}
	return s.wrap(data), nil
}

func (s *EnvironmentService) List(ctx context.Context, params *api.ListEnvironmentsParams) ([]*Environment, error) {
	path, err := s.client.orgPath("/environments")
	if err != nil {
		return nil, err
	}
	if params != nil {
		q := url.Values{}
		if params.Type != nil {
			q.Set("type", *params.Type)
		}
		if params.State != nil {
			q.Set("state", *params.State)
		}
		if qs := q.Encode(); qs != "" {
			path += "?" + qs
		}
	}
	var list []api.Environment
	if err := s.client.do(ctx, "GET", path, nil, &list); err != nil {
		return nil, err
	}
	envs := make([]*Environment, len(list))
	for i, d := range list {
		envs[i] = s.wrap(d)
	}
	return envs, nil
}

func (s *EnvironmentService) Get(ctx context.Context, envID string) (*Environment, error) {
	path, err := s.client.orgPath("/environments/%s", envID)
	if err != nil {
		return nil, err
	}
	var data api.Environment
	if err := s.client.do(ctx, "GET", path, nil, &data); err != nil {
		return nil, err
	}
	return s.wrap(data), nil
}

func (s *EnvironmentService) Delete(ctx context.Context, envID string) error {
	path, err := s.client.orgPath("/environments/%s", envID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *EnvironmentService) Stop(ctx context.Context, envID string) error {
	path, err := s.client.orgPath("/environments/%s/stop", envID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "POST", path, nil, nil)
}

func (s *EnvironmentService) Start(ctx context.Context, envID string) error {
	path, err := s.client.orgPath("/environments/%s/start", envID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "POST", path, nil, nil)
}

func (s *EnvironmentService) Exec(ctx context.Context, envID string, req api.ExecBody) (*api.ExecResult, error) {
	path, err := s.client.orgPath("/environments/%s/exec", envID)
	if err != nil {
		return nil, err
	}
	var result api.ExecResult
	if err := s.client.do(ctx, "POST", path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *EnvironmentService) ListFiles(ctx context.Context, envID string, dirPath *string) ([]api.FileEntry, error) {
	path, err := s.client.orgPath("/environments/%s/files", envID)
	if err != nil {
		return nil, err
	}
	if dirPath != nil {
		path += "?path=" + url.QueryEscape(*dirPath)
	}
	var files []api.FileEntry
	if err := s.client.do(ctx, "GET", path, nil, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (s *EnvironmentService) UploadFile(ctx context.Context, envID string, src io.Reader, remotePath string) error {
	path, err := s.client.orgPath("/environments/%s/files/upload?path=%s", envID, url.QueryEscape(remotePath))
	if err != nil {
		return err
	}
	resp, err := s.client.doRaw(ctx, "POST", path, src, "application/octet-stream")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (s *EnvironmentService) DownloadFile(ctx context.Context, envID string, remotePath string) (io.ReadCloser, error) {
	path, err := s.client.orgPath("/environments/%s/files/download?path=%s", envID, url.QueryEscape(remotePath))
	if err != nil {
		return nil, err
	}
	resp, err := s.client.doRaw(ctx, "GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s *EnvironmentService) ListPorts(ctx context.Context, envID string) ([]api.EnvironmentPort, error) {
	path, err := s.client.orgPath("/environments/%s/ports", envID)
	if err != nil {
		return nil, err
	}
	var ports []api.EnvironmentPort
	if err := s.client.do(ctx, "GET", path, nil, &ports); err != nil {
		return nil, err
	}
	return ports, nil
}

func (s *EnvironmentService) CreatePort(ctx context.Context, envID string, req api.CreatePortBody) (*api.EnvironmentPort, error) {
	path, err := s.client.orgPath("/environments/%s/ports", envID)
	if err != nil {
		return nil, err
	}
	var port api.EnvironmentPort
	if err := s.client.do(ctx, "POST", path, req, &port); err != nil {
		return nil, err
	}
	return &port, nil
}

func (s *EnvironmentService) DeletePort(ctx context.Context, envID, portID string) error {
	path, err := s.client.orgPath("/environments/%s/ports/%s", envID, portID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

func (s *EnvironmentService) wrap(data api.Environment) *Environment {
	orgID, _ := s.client.resolveOrgID()
	return &Environment{
		Environment: data,
		client:      s.client,
		orgID:       orgID,
	}
}

func orgPath(orgID, pattern string, args ...any) string {
	allArgs := make([]any, 0, len(args)+1)
	allArgs = append(allArgs, orgID)
	allArgs = append(allArgs, args...)
	return fmt.Sprintf("/organizations/%s"+pattern, allArgs...)
}
