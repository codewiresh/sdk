import type { HttpClient } from "./http.js";
import type {
  SecretProject,
  SecretMetadata,
  CreateSecretProjectRequest,
  SetProjectSecretRequest,
} from "./types.js";

export class SecretProjects {
  constructor(private http: HttpClient) {}

  create(body: CreateSecretProjectRequest): Promise<SecretProject> {
    return this.http.post(this.http.orgPath("/secret-projects"), body);
  }

  list(): Promise<SecretProject[]> {
    return this.http.get(this.http.orgPath("/secret-projects"));
  }

  delete(projectId: string): Promise<void> {
    return this.http.delete(
      this.http.orgPath(`/secret-projects/${projectId}`),
    );
  }

  async listSecrets(projectId: string): Promise<SecretMetadata[]> {
    const resp = await this.http.get<{ secrets: SecretMetadata[] }>(
      this.http.orgPath(`/secret-projects/${projectId}/secrets`),
    );
    return resp.secrets ?? [];
  }

  setSecret(projectId: string, body: SetProjectSecretRequest): Promise<void> {
    return this.http.put(
      this.http.orgPath(`/secret-projects/${projectId}/secrets`),
      body,
    );
  }

  deleteSecret(projectId: string, key: string): Promise<void> {
    return this.http.delete(
      this.http.orgPath(
        `/secret-projects/${projectId}/secrets/${encodeURIComponent(key)}`,
      ),
    );
  }
}
