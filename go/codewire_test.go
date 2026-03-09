package codewire_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	codewire "codewire.sh/sdk-go"
	"codewire.sh/sdk-go/internal/api"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*codewire.Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := codewire.New("test-key",
		codewire.WithBaseURL(srv.URL),
		codewire.WithOrgID("org_test"),
	)
	return client, srv
}

func TestEnvironmentsCRUD(t *testing.T) {
	envJSON := api.Environment{
		Id: "env_123", OrgId: "org_test", State: "running",
		DesiredState: "running", Type: "sandbox", CreatedBy: "user_1",
		TemplateId: "tmpl_1", Recoverable: true, StateChangedAt: "2026-01-01T00:00:00Z",
		CpuMillicores: 2000, MemoryMb: 4096, DiskGb: 20,
		TotalRunningSeconds: 0, CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing auth header")
		}

		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/environments"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(envJSON)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/environments/env_123"):
			json.NewEncoder(w).Encode(envJSON)

		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/environments"):
			json.NewEncoder(w).Encode([]api.Environment{envJSON})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/environments/env_123"):
			w.WriteHeader(204)

		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/stop"):
			w.WriteHeader(204)

		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/start"):
			w.WriteHeader(204)

		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	// Create
	slug := "python"
	env, err := client.Environments.Create(ctx, api.CreateEnvironmentBody{
		TemplateSlug: &slug,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if env.Id != "env_123" {
		t.Errorf("expected env_123, got %s", env.Id)
	}
	if env.State != "running" {
		t.Errorf("expected running, got %s", env.State)
	}

	// Get
	env2, err := client.Environments.Get(ctx, "env_123")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if env2.Id != "env_123" {
		t.Errorf("expected env_123, got %s", env2.Id)
	}

	// List
	envs, err := client.Environments.List(ctx, nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(envs) != 1 {
		t.Errorf("expected 1 env, got %d", len(envs))
	}

	// Stop/Start
	if err := client.Environments.Stop(ctx, "env_123"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := client.Environments.Start(ctx, "env_123"); err != nil {
		t.Fatalf("start: %v", err)
	}

	// Delete
	if err := client.Environments.Delete(ctx, "env_123"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestEnvironmentExec(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/exec") {
			var body api.ExecBody
			json.NewDecoder(r.Body).Decode(&body)
			if body.Command == nil || len(*body.Command) == 0 {
				t.Error("expected command")
			}
			json.NewEncoder(w).Encode(api.ExecResult{
				ExitCode: 0,
				Stdout:   "hello\n",
				Stderr:   "",
			})
			return
		}
		w.WriteHeader(404)
	})

	ctx := context.Background()
	cmd := []string{"echo", "hello"}
	result, err := client.Environments.Exec(ctx, "env_123", api.ExecBody{
		Command: &cmd,
	})
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit 0, got %d", result.ExitCode)
	}
	if result.Stdout != "hello\n" {
		t.Errorf("unexpected stdout: %q", result.Stdout)
	}
}

func TestEnvironmentFileUploadDownload(t *testing.T) {
	content := "file contents here"
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/files/upload"):
			path := r.URL.Query().Get("path")
			if path != "/workspace/test.txt" {
				t.Errorf("expected path /workspace/test.txt, got %s", path)
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != content {
				t.Errorf("expected %q, got %q", content, string(body))
			}
			w.WriteHeader(200)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/files/download"):
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte(content))

		default:
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	// Upload
	err := client.Environments.UploadFile(ctx, "env_123", strings.NewReader(content), "/workspace/test.txt")
	if err != nil {
		t.Fatalf("upload: %v", err)
	}

	// Download
	rc, err := client.Environments.DownloadFile(ctx, "env_123", "/workspace/test.txt")
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	defer rc.Close()
	got, _ := io.ReadAll(rc)
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
}

func TestEnvironmentWrapper(t *testing.T) {
	envJSON := api.Environment{
		Id: "env_w", OrgId: "org_test", State: "running",
		DesiredState: "running", Type: "sandbox", CreatedBy: "user_1",
		TemplateId: "tmpl_1", Recoverable: true, StateChangedAt: "2026-01-01T00:00:00Z",
		CpuMillicores: 1000, MemoryMb: 2048, DiskGb: 10,
		TotalRunningSeconds: 0, CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/environments"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(envJSON)
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/exec"):
			json.NewEncoder(w).Encode(api.ExecResult{Stdout: "ok", ExitCode: 0})
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/stop"):
			w.WriteHeader(204)
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/environments/env_w"):
			w.WriteHeader(204)
		default:
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()
	slug := "go"
	env, err := client.Environments.Create(ctx, api.CreateEnvironmentBody{TemplateSlug: &slug})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Use wrapper methods
	cmd := []string{"go", "version"}
	result, err := env.Exec(ctx, api.ExecBody{Command: &cmd})
	if err != nil {
		t.Fatalf("exec: %v", err)
	}
	if result.Stdout != "ok" {
		t.Errorf("expected ok, got %s", result.Stdout)
	}

	if err := env.Stop(ctx); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := env.Remove(ctx); err != nil {
		t.Fatalf("remove: %v", err)
	}
}

func TestTemplatesCRUD(t *testing.T) {
	tmplJSON := api.EnvironmentTemplate{
		Id: "tmpl_1", OrgId: "org_test", Name: "Go", Type: "sandbox",
		BuildStatus: "ready", DefaultCpuMillicores: 2000,
		DefaultMemoryMb: 4096, DefaultDiskGb: 20, Official: false,
		CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/templates"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(tmplJSON)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/templates/tmpl_1"):
			json.NewEncoder(w).Encode(tmplJSON)
		case r.Method == "PATCH" && strings.HasSuffix(r.URL.Path, "/templates/tmpl_1"):
			json.NewEncoder(w).Encode(tmplJSON)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/templates"):
			json.NewEncoder(w).Encode([]api.EnvironmentTemplate{tmplJSON})
		case r.Method == "DELETE" && strings.HasSuffix(r.URL.Path, "/templates/tmpl_1"):
			w.WriteHeader(204)
		default:
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	tmpl, err := client.Templates.Create(ctx, api.CreateTemplateBody{Name: "Go", Type: "sandbox"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if tmpl.Name != "Go" {
		t.Errorf("expected Go, got %s", tmpl.Name)
	}

	_, err = client.Templates.Get(ctx, "tmpl_1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	newName := "Golang"
	_, err = client.Templates.Update(ctx, "tmpl_1", api.UpdateTemplateBody{Name: &newName})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	list, err := client.Templates.List(ctx, nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 template, got %d", len(list))
	}

	if err := client.Templates.Delete(ctx, "tmpl_1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestAPIKeys(t *testing.T) {
	keyJSON := api.CreateAPIKeyResponse{
		Id: "key_1", UserId: "user_1", OrgId: "org_test", Name: "test",
		Key: "cw_abc123", KeyPrefix: "cw_abc", CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/api-keys"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(keyJSON)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/api-keys"):
			json.NewEncoder(w).Encode([]api.APIKey{{
				Id: "key_1", UserId: "user_1", OrgId: "org_test",
				Name: "test", KeyPrefix: "cw_abc", CreatedAt: "2026-01-01T00:00:00Z",
			}})
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/api-keys/key_1"):
			w.WriteHeader(204)
		default:
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	resp, err := client.APIKeys.Create(ctx, api.CreateAPIKeyBody{Name: "test"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Key != "cw_abc123" {
		t.Errorf("expected cw_abc123, got %s", resp.Key)
	}

	keys, err := client.APIKeys.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}

	if err := client.APIKeys.Delete(ctx, "key_1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestSecrets(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/secrets":
			if r.URL.Query().Get("org_id") != "org_test" {
				t.Errorf("missing org_id query param")
			}
			json.NewEncoder(w).Encode(api.ListSecretsOutputBody{
				Secrets: &[]api.SecretMetadata{{Key: "DB_URL"}},
			})
		case r.Method == "PUT" && r.URL.Path == "/secrets":
			var body api.SetSecretInputBody
			json.NewDecoder(r.Body).Decode(&body)
			if body.OrgId != "org_test" {
				t.Errorf("expected org_test, got %s", body.OrgId)
			}
			json.NewEncoder(w).Encode(api.SetSecretOutputBody{Status: "ok"})
		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/secrets/"):
			if r.URL.Query().Get("org_id") != "org_test" {
				t.Errorf("missing org_id query param on delete")
			}
			w.WriteHeader(204)
		case r.Method == "GET" && r.URL.Path == "/user/secrets":
			json.NewEncoder(w).Encode(api.ListSecretsOutputBody{
				Secrets: &[]api.SecretMetadata{{Key: "TOKEN"}},
			})
		case r.Method == "PUT" && r.URL.Path == "/user/secrets":
			w.WriteHeader(204)
		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/user/secrets/"):
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	secrets, err := client.Secrets.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(secrets) != 1 || secrets[0].Key != "DB_URL" {
		t.Errorf("unexpected secrets: %v", secrets)
	}

	if err := client.Secrets.Set(ctx, "DB_URL", "postgres://..."); err != nil {
		t.Fatalf("set: %v", err)
	}

	if err := client.Secrets.Delete(ctx, "DB_URL"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	userSecrets, err := client.Secrets.ListUser(ctx)
	if err != nil {
		t.Fatalf("list user: %v", err)
	}
	if len(userSecrets) != 1 {
		t.Errorf("expected 1 user secret, got %d", len(userSecrets))
	}

	if err := client.Secrets.SetUser(ctx, "TOKEN", "abc"); err != nil {
		t.Fatalf("set user: %v", err)
	}

	if err := client.Secrets.DeleteUser(ctx, "TOKEN"); err != nil {
		t.Fatalf("delete user: %v", err)
	}
}

func TestSecretProjects(t *testing.T) {
	projJSON := api.SecretProject{
		Id: "proj_1", OrgId: "org_test", Name: "myapp",
		CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/secret-projects"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(projJSON)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/secret-projects"):
			json.NewEncoder(w).Encode([]api.SecretProject{projJSON})
		case r.Method == "DELETE" && strings.HasSuffix(r.URL.Path, "/secret-projects/proj_1"):
			w.WriteHeader(204)
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/secrets"):
			json.NewEncoder(w).Encode(api.ListProjectSecretsOutputBody{
				Secrets: &[]api.SecretMetadata{{Key: "API_KEY"}},
			})
		case r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/secrets"):
			w.WriteHeader(204)
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/secrets/API_KEY"):
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	proj, err := client.SecretProjects.Create(ctx, api.CreateSecretProjectInputBody{Name: "myapp"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if proj.Name != "myapp" {
		t.Errorf("expected myapp, got %s", proj.Name)
	}

	list, err := client.SecretProjects.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 project, got %d", len(list))
	}

	secrets, err := client.SecretProjects.ListSecrets(ctx, "proj_1")
	if err != nil {
		t.Fatalf("list secrets: %v", err)
	}
	if len(secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(secrets))
	}

	if err := client.SecretProjects.SetSecret(ctx, "proj_1", "API_KEY", "sk_123"); err != nil {
		t.Fatalf("set secret: %v", err)
	}

	if err := client.SecretProjects.DeleteSecret(ctx, "proj_1", "API_KEY"); err != nil {
		t.Fatalf("delete secret: %v", err)
	}

	if err := client.SecretProjects.Delete(ctx, "proj_1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestErrorHandling(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"detail": "environment not found"})
	})

	_, err := client.Environments.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *codewire.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *codewire.Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Detail != "environment not found" {
		t.Errorf("expected 'environment not found', got %q", apiErr.Detail)
	}
}

func TestMissingOrgID(t *testing.T) {
	client := codewire.New("test-key")
	_, err := client.Environments.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for missing orgID")
	}
	if !strings.Contains(err.Error(), "orgID is required") {
		t.Errorf("expected orgID error, got: %v", err)
	}
}

func TestPorts(t *testing.T) {
	portJSON := api.EnvironmentPort{
		Id: "port_1", EnvironmentId: "env_123", Port: 8080,
		Label: "web", Access: "public", CreatedAt: "2026-01-01T00:00:00Z",
	}

	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/ports"):
			json.NewEncoder(w).Encode([]api.EnvironmentPort{portJSON})
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/ports"):
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(portJSON)
		case r.Method == "DELETE" && strings.HasSuffix(r.URL.Path, "/ports/port_1"):
			w.WriteHeader(204)
		default:
			w.WriteHeader(404)
		}
	})

	ctx := context.Background()

	ports, err := client.Environments.ListPorts(ctx, "env_123")
	if err != nil {
		t.Fatalf("list ports: %v", err)
	}
	if len(ports) != 1 || ports[0].Port != 8080 {
		t.Errorf("unexpected ports: %v", ports)
	}

	access := "public"
	port, err := client.Environments.CreatePort(ctx, "env_123", api.CreatePortBody{
		Port: 8080, Label: "web", Access: &access,
	})
	if err != nil {
		t.Fatalf("create port: %v", err)
	}
	if port.Label != "web" {
		t.Errorf("expected web, got %s", port.Label)
	}

	if err := client.Environments.DeletePort(ctx, "env_123", "port_1"); err != nil {
		t.Fatalf("delete port: %v", err)
	}
}
