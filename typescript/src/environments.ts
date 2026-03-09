import type { HttpClient } from "./http.js";
import { Environment } from "./environment.js";
import type {
  CreateEnvironmentRequest,
  EnvironmentData,
  ExecRequest,
  ExecResult,
  FileEntry,
  EnvironmentPort,
  CreatePortRequest,
  ListEnvironmentsParams,
} from "./types.js";

export class Environments {
  constructor(private http: HttpClient) {}

  async create(body: CreateEnvironmentRequest): Promise<Environment> {
    const data = await this.http.post<EnvironmentData>(
      this.http.orgPath("/environments"),
      body,
    );
    return new Environment(data, this.http);
  }

  async list(params?: ListEnvironmentsParams): Promise<Environment[]> {
    const qs = new URLSearchParams();
    if (params?.type) qs.set("type", params.type);
    if (params?.state) qs.set("state", params.state);
    const q = qs.toString();
    const path = this.http.orgPath(`/environments${q ? `?${q}` : ""}`);
    const list = await this.http.get<EnvironmentData[]>(path);
    return list.map((d) => new Environment(d, this.http));
  }

  async get(envId: string): Promise<Environment> {
    const data = await this.http.get<EnvironmentData>(
      this.http.orgPath(`/environments/${envId}`),
    );
    return new Environment(data, this.http);
  }

  delete(envId: string): Promise<void> {
    return this.http.delete(this.http.orgPath(`/environments/${envId}`));
  }

  stop(envId: string): Promise<void> {
    return this.http.post(this.http.orgPath(`/environments/${envId}/stop`));
  }

  start(envId: string): Promise<void> {
    return this.http.post(this.http.orgPath(`/environments/${envId}/start`));
  }

  exec(envId: string, req: ExecRequest): Promise<ExecResult> {
    return this.http.post(
      this.http.orgPath(`/environments/${envId}/exec`),
      req,
    );
  }

  listFiles(envId: string, path?: string): Promise<FileEntry[]> {
    const qs = path ? `?path=${encodeURIComponent(path)}` : "";
    return this.http.get(
      this.http.orgPath(`/environments/${envId}/files${qs}`),
    );
  }

  async uploadFile(
    envId: string,
    src: Uint8Array,
    remotePath: string,
  ): Promise<void> {
    const qs = `?path=${encodeURIComponent(remotePath)}`;
    await this.http.requestRaw(
      "POST",
      this.http.orgPath(`/environments/${envId}/files/upload${qs}`),
      src,
      "application/octet-stream",
    );
  }

  async downloadFile(envId: string, remotePath: string): Promise<ArrayBuffer> {
    const qs = `?path=${encodeURIComponent(remotePath)}`;
    const res = await this.http.requestRaw(
      "GET",
      this.http.orgPath(`/environments/${envId}/files/download${qs}`),
    );
    return res.arrayBuffer();
  }

  listPorts(envId: string): Promise<EnvironmentPort[]> {
    return this.http.get(this.http.orgPath(`/environments/${envId}/ports`));
  }

  createPort(envId: string, req: CreatePortRequest): Promise<EnvironmentPort> {
    return this.http.post(
      this.http.orgPath(`/environments/${envId}/ports`),
      req,
    );
  }

  deletePort(envId: string, portId: string): Promise<void> {
    return this.http.delete(
      this.http.orgPath(`/environments/${envId}/ports/${portId}`),
    );
  }
}
