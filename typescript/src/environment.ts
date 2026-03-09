import type { HttpClient } from "./http.js";
import type {
  EnvironmentData,
  ExecRequest,
  ExecResult,
  FileEntry,
  EnvironmentPort,
  CreatePortRequest,
} from "./types.js";

export class Environment implements EnvironmentData {
  readonly id: string;
  readonly org_id: string;
  readonly created_by: string;
  readonly template_id: string;
  readonly type: string;
  readonly name: string | undefined;
  readonly state: string;
  readonly desired_state: string;
  readonly error_reason: string | undefined;
  readonly recoverable: boolean;
  readonly state_changed_at: string;
  readonly retry_count: number;
  readonly ttl_seconds: number | undefined;
  readonly shutdown_at: string | undefined;
  readonly cpu_millicores: number;
  readonly memory_mb: number;
  readonly disk_gb: number;
  readonly total_running_seconds: number;
  readonly created_at: string;
  readonly started_at: string | undefined;
  readonly stopped_at: string | undefined;
  readonly destroyed_at: string | undefined;
  readonly ports: EnvironmentPort[] | undefined;
  readonly template_name: string | undefined;

  private http: HttpClient;
  private _orgId: string;

  constructor(data: EnvironmentData, http: HttpClient) {
    Object.assign(this, data);
    this.id = data.id;
    this.org_id = data.org_id;
    this.created_by = data.created_by;
    this.template_id = data.template_id;
    this.type = data.type;
    this.state = data.state;
    this.desired_state = data.desired_state;
    this.recoverable = data.recoverable;
    this.state_changed_at = data.state_changed_at;
    this.retry_count = data.retry_count;
    this.cpu_millicores = data.cpu_millicores;
    this.memory_mb = data.memory_mb;
    this.disk_gb = data.disk_gb;
    this.total_running_seconds = data.total_running_seconds;
    this.created_at = data.created_at;
    this.http = http;
    this._orgId = http.orgId;
  }

  private path(suffix: string): string {
    return `/organizations/${this._orgId}/environments/${this.id}${suffix}`;
  }

  async exec(req: ExecRequest): Promise<ExecResult> {
    return this.http.post(this.path("/exec"), req);
  }

  async start(): Promise<void> {
    await this.http.post(this.path("/start"));
  }

  async stop(): Promise<void> {
    await this.http.post(this.path("/stop"));
  }

  async remove(): Promise<void> {
    await this.http.delete(this.path(""));
  }

  async upload(src: Uint8Array, remotePath: string): Promise<void> {
    const qs = `?path=${encodeURIComponent(remotePath)}`;
    await this.http.requestRaw(
      "POST",
      this.path(`/files/upload${qs}`),
      src,
      "application/octet-stream",
    );
  }

  async download(remotePath: string): Promise<ArrayBuffer> {
    const qs = `?path=${encodeURIComponent(remotePath)}`;
    const res = await this.http.requestRaw("GET", this.path(`/files/download${qs}`));
    return res.arrayBuffer();
  }

  async listFiles(dirPath?: string): Promise<FileEntry[]> {
    const qs = dirPath ? `?path=${encodeURIComponent(dirPath)}` : "";
    return this.http.get(this.path(`/files${qs}`));
  }

  async listPorts(): Promise<EnvironmentPort[]> {
    return this.http.get(this.path("/ports"));
  }

  async createPort(req: CreatePortRequest): Promise<EnvironmentPort> {
    return this.http.post(this.path("/ports"), req);
  }

  async deletePort(portId: string): Promise<void> {
    await this.http.delete(this.path(`/ports/${portId}`));
  }
}
