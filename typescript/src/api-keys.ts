import type { HttpClient } from "./http.js";
import type {
  APIKeyData,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
} from "./types.js";

export class APIKeys {
  constructor(private http: HttpClient) {}

  create(body: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
    return this.http.post(this.http.orgPath("/api-keys"), body);
  }

  list(): Promise<APIKeyData[]> {
    return this.http.get(this.http.orgPath("/api-keys"));
  }

  delete(keyId: string): Promise<void> {
    return this.http.delete(this.http.orgPath(`/api-keys/${keyId}`));
  }
}
