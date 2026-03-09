package codewire

import (
	"context"
	"io"
	"net/url"

	"codewire.sh/sdk-go/internal/api"
)

// Environment wraps the generated api.Environment data with convenience methods.
type Environment struct {
	api.Environment
	client *Client
	orgID  string
}

func (e *Environment) Exec(ctx context.Context, req api.ExecBody) (*api.ExecResult, error) {
	var result api.ExecResult
	err := e.client.do(ctx, "POST", orgPath(e.orgID, "/environments/%s/exec", e.Id), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (e *Environment) Start(ctx context.Context) error {
	return e.client.do(ctx, "POST", orgPath(e.orgID, "/environments/%s/start", e.Id), nil, nil)
}

func (e *Environment) Stop(ctx context.Context) error {
	return e.client.do(ctx, "POST", orgPath(e.orgID, "/environments/%s/stop", e.Id), nil, nil)
}

func (e *Environment) Remove(ctx context.Context) error {
	return e.client.do(ctx, "DELETE", orgPath(e.orgID, "/environments/%s", e.Id), nil, nil)
}

func (e *Environment) Upload(ctx context.Context, src io.Reader, remotePath string) error {
	path := orgPath(e.orgID, "/environments/%s/files/upload?path=%s", e.Id, url.QueryEscape(remotePath))
	resp, err := e.client.doRaw(ctx, "POST", path, src, "application/octet-stream")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (e *Environment) Download(ctx context.Context, remotePath string) (io.ReadCloser, error) {
	path := orgPath(e.orgID, "/environments/%s/files/download?path=%s", e.Id, url.QueryEscape(remotePath))
	resp, err := e.client.doRaw(ctx, "GET", path, nil, "")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (e *Environment) ListFiles(ctx context.Context, dirPath *string) ([]api.FileEntry, error) {
	path := orgPath(e.orgID, "/environments/%s/files", e.Id)
	if dirPath != nil {
		path += "?path=" + url.QueryEscape(*dirPath)
	}
	var files []api.FileEntry
	if err := e.client.do(ctx, "GET", path, nil, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (e *Environment) ListPorts(ctx context.Context) ([]api.EnvironmentPort, error) {
	var ports []api.EnvironmentPort
	err := e.client.do(ctx, "GET", orgPath(e.orgID, "/environments/%s/ports", e.Id), nil, &ports)
	if err != nil {
		return nil, err
	}
	return ports, nil
}

func (e *Environment) CreatePort(ctx context.Context, req api.CreatePortBody) (*api.EnvironmentPort, error) {
	var port api.EnvironmentPort
	err := e.client.do(ctx, "POST", orgPath(e.orgID, "/environments/%s/ports", e.Id), req, &port)
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func (e *Environment) DeletePort(ctx context.Context, portID string) error {
	return e.client.do(ctx, "DELETE", orgPath(e.orgID, "/environments/%s/ports/%s", e.Id, portID), nil, nil)
}
