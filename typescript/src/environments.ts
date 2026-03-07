import type { HttpClient } from "./http.js";
import type {
  CreateEnvironmentRequest,
  EnvironmentData,
  ListEnvironmentsParams,
  StatusResponse,
} from "./types.js";

export class Environments {
  constructor(private http: HttpClient) {}

  create(
    orgId: string,
    body: CreateEnvironmentRequest,
  ): Promise<EnvironmentData> {
    return this.http.post(`/organizations/${orgId}/environments`, body);
  }

  list(
    orgId: string,
    params?: ListEnvironmentsParams,
  ): Promise<EnvironmentData[]> {
    const qs = new URLSearchParams();
    if (params?.type) qs.set("type", params.type);
    if (params?.state) qs.set("state", params.state);
    const q = qs.toString();
    return this.http.get(
      `/organizations/${orgId}/environments${q ? `?${q}` : ""}`,
    );
  }

  get(orgId: string, envId: string): Promise<EnvironmentData> {
    return this.http.get(`/organizations/${orgId}/environments/${envId}`);
  }

  delete(orgId: string, envId: string): Promise<void> {
    return this.http.delete(`/organizations/${orgId}/environments/${envId}`);
  }

  stop(orgId: string, envId: string): Promise<StatusResponse> {
    return this.http.post(
      `/organizations/${orgId}/environments/${envId}/stop`,
    );
  }

  start(orgId: string, envId: string): Promise<StatusResponse> {
    return this.http.post(
      `/organizations/${orgId}/environments/${envId}/start`,
    );
  }
}
