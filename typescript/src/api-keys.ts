import type { HttpClient } from "./http.js";
import type {
  APIKeyData,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
} from "./types.js";

export class APIKeys {
  constructor(private http: HttpClient) {}

  create(
    orgId: string,
    body: CreateAPIKeyRequest,
  ): Promise<CreateAPIKeyResponse> {
    return this.http.post(`/organizations/${orgId}/api-keys`, body);
  }

  list(orgId: string): Promise<APIKeyData[]> {
    return this.http.get(`/organizations/${orgId}/api-keys`);
  }

  delete(orgId: string, keyId: string): Promise<void> {
    return this.http.delete(`/organizations/${orgId}/api-keys/${keyId}`);
  }
}
