import { HttpClient } from "./http.js";
import { Environments } from "./environments.js";
import { Templates } from "./templates.js";
import { APIKeys } from "./api-keys.js";
import type { ClientOptions } from "./types.js";

export class Codewire {
  readonly environments: Environments;
  readonly templates: Templates;
  readonly apiKeys: APIKeys;

  constructor(opts: ClientOptions = {}) {
    const apiKey = opts.apiKey ?? process.env.CODEWIRE_API_KEY;
    if (!apiKey) {
      throw new Error(
        "Codewire API key is required. Pass apiKey or set CODEWIRE_API_KEY.",
      );
    }
    const baseUrl = (opts.baseUrl ?? "https://api.codewire.sh").replace(
      /\/$/,
      "",
    );
    const http = new HttpClient(baseUrl, apiKey);
    this.environments = new Environments(http);
    this.templates = new Templates(http);
    this.apiKeys = new APIKeys(http);
  }
}

export { Codewire as default };
export { CodewireError } from "./types.js";
export type {
  ClientOptions,
  CreateEnvironmentRequest,
  EnvironmentData,
  TemplateData,
  CreateTemplateRequest,
  APIKeyData,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  ListEnvironmentsParams,
  ListTemplatesParams,
  StatusResponse,
} from "./types.js";
