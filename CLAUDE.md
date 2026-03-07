# Codewire SDK

Official client libraries for the Codewire API.

## Structure

- `go/` — Go SDK (`codewire.sh/sdk-go`)
- `python/` — Python SDK (`codewire`)
- `typescript/` — TypeScript SDK (`@codewire/sdk`)

## Repos

- **GitHub**: `codewiresh/sdk`
- **Gitea**: `codewire/sdk` (git.noel.sh)

## API Coverage

All three SDKs implement the same surface:

- **Environments**: create, list, get, delete, stop, start
- **Templates**: create, list, delete
- **API Keys**: create, list, delete

## Auth

SDKs authenticate with `cw_` prefixed API keys via Bearer token.
Keys are created at `POST /organizations/{orgId}/api-keys`.

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
