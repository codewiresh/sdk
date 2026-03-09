import type { components } from "./types.gen.js";

// Re-export generated schema types with friendly names
export type EnvironmentData = components["schemas"]["Environment"];
export type TemplateData = components["schemas"]["EnvironmentTemplate"];
export type APIKeyData = components["schemas"]["APIKey"];
export type CreateEnvironmentRequest = components["schemas"]["CreateEnvironmentBody"];
export type CreateTemplateRequest = components["schemas"]["CreateTemplateBody"];
export type UpdateTemplateRequest = components["schemas"]["UpdateTemplateBody"];
export type CreateAPIKeyRequest = components["schemas"]["CreateAPIKeyBody"];
export type CreateAPIKeyResponse = components["schemas"]["CreateAPIKeyResponse"];
export type ExecRequest = components["schemas"]["ExecBody"];
export type ExecResult = components["schemas"]["ExecResult"];
export type FileEntry = components["schemas"]["FileEntry"];
export type EnvironmentPort = components["schemas"]["EnvironmentPort"];
export type CreatePortRequest = components["schemas"]["CreatePortBody"];
export type SecretMetadata = components["schemas"]["SecretMetadata"];
export type SecretProject = components["schemas"]["SecretProject"];
export type CreateSecretProjectRequest = components["schemas"]["CreateSecretProjectInputBody"];
export type SetSecretRequest = components["schemas"]["SetSecretInputBody"];
export type SetUserSecretRequest = components["schemas"]["SetUserSecretInputBody"];
export type SetProjectSecretRequest = components["schemas"]["SetProjectSecretInputBody"];
export type StatusResponse = components["schemas"]["StatusResponse"];

export interface ClientOptions {
  apiKey?: string;
  baseUrl?: string;
  orgId?: string;
}

export interface ListEnvironmentsParams {
  type?: string;
  state?: string;
}

export interface ListTemplatesParams {
  type?: string;
}

export class CodewireError extends Error {
  constructor(
    public status: number,
    public detail: string,
  ) {
    super(`Codewire API error ${status}: ${detail}`);
    this.name = "CodewireError";
  }
}
