# Codewire SDK

Official client libraries for the Codewire API.

## Structure

- `go/` — Go SDK (`codewire.sh/sdk-go`)
- `python/` — Python SDK (`codewire`)
- `typescript/` — TypeScript SDK (`@codewire/sdk`)
- `scripts/filter-openapi.py` — Filters server OpenAPI spec for SDK-relevant endpoints

## Repos

- **GitHub**: `codewiresh/sdk`
- **Gitea**: `codewire/sdk` (git.noel.sh)

## Codegen

Types are generated from the server's OpenAPI spec. Run from `sdk/`:

```bash
make generate           # Filter OpenAPI + generate Go/TS/Python types
make generate-go        # Go only (oapi-codegen)
make generate-ts        # TypeScript only (openapi-typescript)
make generate-py        # Python only (datamodel-code-generator)
make check              # Build/typecheck all languages
```

Generated files (gitignored):
- `openapi-sdk.json` — Filtered OpenAPI spec
- `go/internal/api/types_gen.go` — Go types
- `typescript/src/types.gen.ts` — TypeScript types
- `python/codewire/models.py` — Pydantic models

## API Coverage

All three SDKs implement the same surface:

- **Environments**: create, list, get, delete, start, stop, exec, list_files, upload_file, download_file, list_ports, create_port, delete_port
- **Templates**: create, list, get, update, delete
- **API Keys**: create, list, delete
- **Secrets**: list, set, delete (org-scoped); list_user, set_user, delete_user
- **Secret Projects**: create, list, delete, list_secrets, set_secret, delete_secret

## Architecture

- **Generated types**: Data shapes from OpenAPI (Go: `internal/api/`, TS: `types.gen.ts`, Python: `models.py`)
- **Hand-written types**: Re-export layer with friendly names (TS: `types.ts`)
- **Services**: Stateless service classes (Environments, Templates, APIKeys, Secrets, SecretProjects)
- **Environment wrapper**: `create()` and `get()` return an Environment object with methods (`exec()`, `start()`, `remove()`, etc.)
- **Implicit orgId**: Set once on client constructor, resolved automatically for all calls

## Auth

SDKs authenticate with `cw_` prefixed API keys via Bearer token.
orgId is required — set via constructor or `CODEWIRE_ORG_ID` env var.

## API Base URL

Default: `https://api.codewire.sh`
Local dev: `http://localhost:8080/api/v1`

## Development

```bash
# TypeScript
cd typescript && npm install && npm run typecheck

# Go
cd go && go build ./...

# Python
cd python && pip install -e .
```

## Notes

- File upload/download endpoints are raw chi handlers, NOT in the OpenAPI spec — they are hand-written in each SDK
- Secrets endpoints use `?org_id=` query parameter (not path param)
- Secret projects use standard path params like other org resources
