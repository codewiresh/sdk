# Codewire SDK

Official client libraries for the [Codewire](https://codewire.sh) API.

| Language | Package | Status |
|----------|---------|--------|
| TypeScript | [`@codewire/sdk`](./typescript/) | Alpha |
| Go | [`codewire.sh/sdk-go`](./go/) | Alpha |
| Python | [`codewire`](./python/) | Alpha |

## Quick start

### TypeScript

```ts
import { Codewire } from "@codewire/sdk";

const cw = new Codewire({ apiKey: process.env.CODEWIRE_API_KEY });

const env = await cw.environments.create("org_xxx", {
  template_id: "tmpl_xxx",
});
console.log(env.id, env.state);
```

### Go

```go
import "codewire.sh/sdk-go"

client := codewire.New(os.Getenv("CODEWIRE_API_KEY"))

env, err := client.Environments.Create(ctx, "org_xxx", codewire.CreateEnvironmentRequest{
    TemplateID: "tmpl_xxx",
})
```

### Python

```python
from codewire import Codewire

cw = Codewire(api_key=os.environ["CODEWIRE_API_KEY"])

env = cw.environments.create("org_xxx", template_id="tmpl_xxx")
print(env.id, env.state)
```

## Authentication

All SDKs authenticate with a Codewire API key. Create one at **Settings > API Keys** in the Codewire dashboard, or via the API:

```
POST /organizations/{orgId}/api-keys
```

Pass the key as `apiKey` when constructing the client. The SDKs also read `CODEWIRE_API_KEY` from the environment.

## API Reference

See the [Codewire API docs](https://codewire.sh/docs/api) for the full endpoint reference.

## License

MIT
