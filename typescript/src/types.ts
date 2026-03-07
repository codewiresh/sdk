export interface ClientOptions {
  apiKey?: string;
  baseUrl?: string;
  orgId?: string;
}

export interface CreateEnvironmentRequest {
  template_id: string;
  name?: string;
  cpu_millicores?: number;
  memory_mb?: number;
  disk_gb?: number;
  ttl_seconds?: number;
  resource_id?: string;
}

export interface EnvironmentData {
  id: string;
  org_id: string;
  created_by: string;
  template_id: string;
  type: "coder" | "sandbox";
  name?: string;
  state: string;
  desired_state: string;
  error_reason?: string;
  recoverable: boolean;
  state_changed_at: string;
  retry_count: number;
  ttl_seconds?: number;
  shutdown_at?: string;
  cpu_millicores: number;
  memory_mb: number;
  disk_gb: number;
  total_running_seconds: number;
  created_at: string;
  started_at?: string;
  stopped_at?: string;
  destroyed_at?: string;
}

export interface TemplateData {
  id: string;
  org_id: string;
  type: "coder" | "sandbox";
  name: string;
  description?: string;
  build_status: string;
  build_error?: string;
  default_cpu_millicores: number;
  default_memory_mb: number;
  default_disk_gb: number;
  default_ttl_seconds?: number;
  created_at: string;
}

export interface CreateTemplateRequest {
  type: "coder" | "sandbox";
  name: string;
  description?: string;
  default_cpu_millicores?: number;
  default_memory_mb?: number;
  default_disk_gb?: number;
  default_ttl_seconds?: number;
}

export interface APIKeyData {
  id: string;
  user_id: string;
  org_id: string;
  name: string;
  key_prefix: string;
  scopes: string[];
  last_used_at?: string;
  expires_at?: string;
  created_at: string;
}

export interface CreateAPIKeyRequest {
  name: string;
  scopes?: string[];
  expires_in_days?: number;
}

export interface CreateAPIKeyResponse extends APIKeyData {
  key: string;
}

export interface ListEnvironmentsParams {
  type?: "coder" | "sandbox";
  state?: string;
}

export interface ListTemplatesParams {
  type?: "coder" | "sandbox";
}

export interface StatusResponse {
  status: string;
}

export class CodewireError extends Error {
  constructor(
    public status: number,
    public detail: string,
  ) {
    super(`Codewire API error ${status}: ${detail}`);
    this.name = "CodewireError";
  }
}
