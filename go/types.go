package codewire

type Environment struct {
	ID                   string  `json:"id"`
	OrgID                string  `json:"org_id"`
	CreatedBy            string  `json:"created_by"`
	TemplateID           string  `json:"template_id"`
	Type                 string  `json:"type"`
	Name                 *string `json:"name,omitempty"`
	State                string  `json:"state"`
	DesiredState         string  `json:"desired_state"`
	ErrorReason          *string `json:"error_reason,omitempty"`
	Recoverable          bool    `json:"recoverable"`
	StateChangedAt       string  `json:"state_changed_at"`
	RetryCount           int     `json:"retry_count"`
	TTLSeconds           *int    `json:"ttl_seconds,omitempty"`
	ShutdownAt           *string `json:"shutdown_at,omitempty"`
	CPUMillicores        int     `json:"cpu_millicores"`
	MemoryMB             int     `json:"memory_mb"`
	DiskGB               int     `json:"disk_gb"`
	TotalRunningSeconds  int     `json:"total_running_seconds"`
	CreatedAt            string  `json:"created_at"`
	StartedAt            *string `json:"started_at,omitempty"`
	StoppedAt            *string `json:"stopped_at,omitempty"`
	DestroyedAt          *string `json:"destroyed_at,omitempty"`
}

type CreateEnvironmentRequest struct {
	TemplateID    string `json:"template_id"`
	Name          string `json:"name,omitempty"`
	CPUMillicores *int   `json:"cpu_millicores,omitempty"`
	MemoryMB      *int   `json:"memory_mb,omitempty"`
	DiskGB        *int   `json:"disk_gb,omitempty"`
	TTLSeconds    *int   `json:"ttl_seconds,omitempty"`
	ResourceID    string `json:"resource_id,omitempty"`
}

type ListEnvironmentsParams struct {
	Type  string
	State string
}

type Template struct {
	ID                   string  `json:"id"`
	OrgID                string  `json:"org_id"`
	Type                 string  `json:"type"`
	Name                 string  `json:"name"`
	Description          *string `json:"description,omitempty"`
	BuildStatus          string  `json:"build_status"`
	BuildError           *string `json:"build_error,omitempty"`
	DefaultCPUMillicores int     `json:"default_cpu_millicores"`
	DefaultMemoryMB      int     `json:"default_memory_mb"`
	DefaultDiskGB        int     `json:"default_disk_gb"`
	DefaultTTLSeconds    *int    `json:"default_ttl_seconds,omitempty"`
	CreatedAt            string  `json:"created_at"`
}

type CreateTemplateRequest struct {
	Type                 string `json:"type"`
	Name                 string `json:"name"`
	Description          string `json:"description,omitempty"`
	DefaultCPUMillicores *int   `json:"default_cpu_millicores,omitempty"`
	DefaultMemoryMB      *int   `json:"default_memory_mb,omitempty"`
	DefaultDiskGB        *int   `json:"default_disk_gb,omitempty"`
	DefaultTTLSeconds    *int   `json:"default_ttl_seconds,omitempty"`
}

type ListTemplatesParams struct {
	Type string
}

type APIKey struct {
	ID         string   `json:"id"`
	UserID     string   `json:"user_id"`
	OrgID      string   `json:"org_id"`
	Name       string   `json:"name"`
	KeyPrefix  string   `json:"key_prefix"`
	Scopes     []string `json:"scopes"`
	LastUsedAt *string  `json:"last_used_at,omitempty"`
	ExpiresAt  *string  `json:"expires_at,omitempty"`
	CreatedAt  string   `json:"created_at"`
}

type CreateAPIKeyRequest struct {
	Name          string   `json:"name"`
	Scopes        []string `json:"scopes,omitempty"`
	ExpiresInDays *int     `json:"expires_in_days,omitempty"`
}

type CreateAPIKeyResponse struct {
	APIKey
	Key string `json:"key"`
}

type StatusResponse struct {
	Status string `json:"status"`
}
