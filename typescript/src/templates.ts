import type { HttpClient } from "./http.js";
import type {
  CreateTemplateRequest,
  ListTemplatesParams,
  TemplateData,
} from "./types.js";

export class Templates {
  constructor(private http: HttpClient) {}

  create(orgId: string, body: CreateTemplateRequest): Promise<TemplateData> {
    return this.http.post(`/organizations/${orgId}/templates`, body);
  }

  list(orgId: string, params?: ListTemplatesParams): Promise<TemplateData[]> {
    const qs = new URLSearchParams();
    if (params?.type) qs.set("type", params.type);
    const q = qs.toString();
    return this.http.get(
      `/organizations/${orgId}/templates${q ? `?${q}` : ""}`,
    );
  }

  delete(orgId: string, templateId: string): Promise<void> {
    return this.http.delete(
      `/organizations/${orgId}/templates/${templateId}`,
    );
  }
}
