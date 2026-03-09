import type { HttpClient } from "./http.js";
import type {
  CreateTemplateRequest,
  UpdateTemplateRequest,
  ListTemplatesParams,
  TemplateData,
} from "./types.js";

export class Templates {
  constructor(private http: HttpClient) {}

  create(body: CreateTemplateRequest): Promise<TemplateData> {
    return this.http.post(this.http.orgPath("/templates"), body);
  }

  list(params?: ListTemplatesParams): Promise<TemplateData[]> {
    const qs = new URLSearchParams();
    if (params?.type) qs.set("type", params.type);
    const q = qs.toString();
    return this.http.get(this.http.orgPath(`/templates${q ? `?${q}` : ""}`));
  }

  get(templateId: string): Promise<TemplateData> {
    return this.http.get(this.http.orgPath(`/templates/${templateId}`));
  }

  update(templateId: string, body: UpdateTemplateRequest): Promise<TemplateData> {
    return this.http.patch(this.http.orgPath(`/templates/${templateId}`), body);
  }

  delete(templateId: string): Promise<void> {
    return this.http.delete(this.http.orgPath(`/templates/${templateId}`));
  }
}
