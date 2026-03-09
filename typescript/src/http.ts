import { CodewireError } from "./types.js";

export class HttpClient {
  readonly orgId: string;

  constructor(
    private baseUrl: string,
    private apiKey: string,
    orgId: string,
  ) {
    this.orgId = orgId;
  }

  resolveOrgId(): string {
    if (!this.orgId) {
      throw new Error(
        "orgId is required — pass orgId in ClientOptions or set CODEWIRE_ORG_ID",
      );
    }
    return this.orgId;
  }

  orgPath(pattern: string): string {
    return `/organizations/${this.resolveOrgId()}${pattern}`;
  }

  private headers(): Record<string, string> {
    return {
      Authorization: `Bearer ${this.apiKey}`,
      "Content-Type": "application/json",
    };
  }

  async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const url = `${this.baseUrl}${path}`;
    const res = await fetch(url, {
      method,
      headers: this.headers(),
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      let detail = res.statusText;
      try {
        const err = await res.json();
        detail = err.detail ?? err.message ?? detail;
      } catch {
        // ignore parse errors
      }
      throw new CodewireError(res.status, detail);
    }

    if (res.status === 204) return undefined as T;
    return res.json() as Promise<T>;
  }

  async requestRaw(
    method: string,
    path: string,
    body?: unknown,
    contentType?: string,
  ): Promise<Response> {
    const url = `${this.baseUrl}${path}`;
    const headers: Record<string, string> = {
      Authorization: `Bearer ${this.apiKey}`,
    };
    if (contentType) headers["Content-Type"] = contentType;

    const res = await fetch(url, {
      method,
      headers,
      body: body as BodyInit | undefined,
    });
    if (!res.ok) {
      let detail = res.statusText;
      try {
        const err = await res.json();
        detail = err.detail ?? err.message ?? detail;
      } catch {
        // ignore
      }
      throw new CodewireError(res.status, detail);
    }
    return res;
  }

  get<T>(path: string): Promise<T> {
    return this.request("GET", path);
  }

  post<T>(path: string, body?: unknown): Promise<T> {
    return this.request("POST", path, body);
  }

  put<T>(path: string, body?: unknown): Promise<T> {
    return this.request("PUT", path, body);
  }

  patch<T>(path: string, body?: unknown): Promise<T> {
    return this.request("PATCH", path, body);
  }

  delete<T>(path: string): Promise<T> {
    return this.request("DELETE", path);
  }
}
