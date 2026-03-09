import { HttpClient } from "./http.js";
import { Environments } from "./environments.js";
import { Templates } from "./templates.js";
import { APIKeys } from "./api-keys.js";
import { Secrets } from "./secrets.js";
import { SecretProjects } from "./secret-projects.js";
import type { ClientOptions } from "./types.js";

export class Codewire {
  readonly environments: Environments;
  readonly templates: Templates;
  readonly apiKeys: APIKeys;
  readonly secrets: Secrets;
  readonly secretProjects: SecretProjects;

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
    const orgId = opts.orgId ?? process.env.CODEWIRE_ORG_ID ?? "";
    const http = new HttpClient(baseUrl, apiKey, orgId);
    this.environments = new Environments(http);
    this.templates = new Templates(http);
    this.apiKeys = new APIKeys(http);
    this.secrets = new Secrets(http);
    this.secretProjects = new SecretProjects(http);
  }
}

export { Codewire as default };
export { CodewireError } from "./types.js";
export { Environment } from "./environment.js";
export type {
  ClientOptions,
  CreateEnvironmentRequest,
  EnvironmentData,
  TemplateData,
  CreateTemplateRequest,
  UpdateTemplateRequest,
  APIKeyData,
  CreateAPIKeyRequest,
  CreateAPIKeyResponse,
  ExecRequest,
  ExecResult,
  FileEntry,
  EnvironmentPort,
  CreatePortRequest,
  SecretMetadata,
  SecretProject,
  CreateSecretProjectRequest,
  SetSecretRequest,
  SetUserSecretRequest,
  SetProjectSecretRequest,
  ListEnvironmentsParams,
  ListTemplatesParams,
  StatusResponse,
} from "./types.js";
