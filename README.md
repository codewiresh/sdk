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

const cw = new Codewire({ apiKey: "cw_...", orgId: "org_xxx" });

// Create and interact with an environment
const env = await cw.environments.create({ template_slug: "python" });
const result = await env.exec({ command: ["python3", "-c", "print('hello')"] });
console.log(result.stdout);
await env.remove();
```

### Go

```go
import codewire "codewire.sh/sdk-go"

client := codewire.New("", codewire.WithOrgID("org_xxx"))

env, _ := client.Environments.Create(ctx, api.CreateEnvironmentBody{
    TemplateSlug: ptr("python"),
})
result, _ := env.Exec(ctx, api.ExecBody{
    Command: &[]string{"python3", "-c", "print('hello')"},
})
fmt.Println(result.Stdout)
env.Remove(ctx)
```

### Python

```python
from codewire import Codewire

cw = Codewire(api_key="cw_...", org_id="org_xxx")

env = cw.environments.create(template_slug="python")
result = env.exec(command=["python3", "-c", "print('hello')"])
print(result["stdout"])
env.remove()
```

## Services

All SDKs provide these services:

- **Environments**: create, list, get, delete, start, stop, exec, files, ports
- **Templates**: create, list, get, update, delete
- **API Keys**: create, list, delete
- **Secrets**: list, set, delete (org + user scopes)
- **Secret Projects**: create, list, delete, manage project secrets

## Environment-as-object

`create()` and `get()` return an `Environment` wrapper with methods:

```
env.exec(...)       # Execute a command
env.start()         # Start the environment
env.stop()          # Stop the environment
env.remove()        # Delete the environment
env.upload(...)     # Upload a file
env.download(...)   # Download a file
env.list_files()    # List files
env.list_ports()    # List exposed ports
```

## Authentication

Pass `apiKey` when constructing the client, or set `CODEWIRE_API_KEY`.
Pass `orgId` when constructing the client, or set `CODEWIRE_ORG_ID`.

## Codegen

Types are generated from the server's OpenAPI spec:

```bash
cd sdk && make generate  # Regenerate all types
cd sdk && make check     # Verify everything compiles
```

## License

MIT
