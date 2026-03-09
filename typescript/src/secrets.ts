import type { HttpClient } from "./http.js";
import type { SecretMetadata } from "./types.js";

export class Secrets {
  constructor(private http: HttpClient) {}

  async list(): Promise<SecretMetadata[]> {
    const orgId = this.http.resolveOrgId();
    const resp = await this.http.get<{ secrets: SecretMetadata[] }>(
      `/secrets?org_id=${encodeURIComponent(orgId)}`,
    );
    return resp.secrets ?? [];
  }

  async set(key: string, value: string): Promise<void> {
    const orgId = this.http.resolveOrgId();
    await this.http.put("/secrets", { org_id: orgId, key, value });
  }

  async delete(key: string): Promise<void> {
    const orgId = this.http.resolveOrgId();
    await this.http.delete(
      `/secrets/${encodeURIComponent(key)}?org_id=${encodeURIComponent(orgId)}`,
    );
  }

  async listUser(): Promise<SecretMetadata[]> {
    const resp = await this.http.get<{ secrets: SecretMetadata[] }>(
      "/user/secrets",
    );
    return resp.secrets ?? [];
  }

  async setUser(key: string, value: string): Promise<void> {
    await this.http.put("/user/secrets", { key, value });
  }

  async deleteUser(key: string): Promise<void> {
    await this.http.delete(`/user/secrets/${encodeURIComponent(key)}`);
  }
}
