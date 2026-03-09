import { createServer, type IncomingMessage, type ServerResponse } from "node:http";
import { Codewire } from "./index.js";
import { Environment } from "./environment.js";
import { CodewireError } from "./types.js";

let passed = 0;
let failed = 0;

function assert(condition: boolean, msg: string) {
  if (!condition) {
    failed++;
    console.error(`  FAIL: ${msg}`);
  } else {
    passed++;
  }
}

function readBody(req: IncomingMessage): Promise<string> {
  return new Promise((resolve) => {
    let data = "";
    req.on("data", (chunk: Buffer) => (data += chunk.toString()));
    req.on("end", () => resolve(data));
  });
}

async function runTests() {
  // Mock server
  const server = createServer(
    async (req: IncomingMessage, res: ServerResponse) => {
      const url = new URL(req.url!, `http://${req.headers.host}`);
      const path = url.pathname;
      const method = req.method!;

      // Auth check
      const auth = req.headers.authorization;
      if (auth !== "Bearer test-key") {
        res.writeHead(401);
        res.end(JSON.stringify({ detail: "unauthorized" }));
        return;
      }

      // Environments
      if (
        method === "POST" &&
        path === "/organizations/org_test/environments"
      ) {
        res.writeHead(201, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            id: "env_123",
            org_id: "org_test",
            state: "running",
            desired_state: "running",
            type: "sandbox",
            created_by: "user_1",
            template_id: "tmpl_1",
            recoverable: true,
            state_changed_at: "2026-01-01T00:00:00Z",
            cpu_millicores: 2000,
            memory_mb: 4096,
            disk_gb: 20,
            total_running_seconds: 0,
            created_at: "2026-01-01T00:00:00Z",
            retry_count: 0,
          }),
        );
        return;
      }

      if (
        method === "GET" &&
        path === "/organizations/org_test/environments/env_123"
      ) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            id: "env_123",
            org_id: "org_test",
            state: "running",
            desired_state: "running",
            type: "sandbox",
            created_by: "user_1",
            template_id: "tmpl_1",
            recoverable: true,
            state_changed_at: "2026-01-01T00:00:00Z",
            cpu_millicores: 2000,
            memory_mb: 4096,
            disk_gb: 20,
            total_running_seconds: 0,
            created_at: "2026-01-01T00:00:00Z",
            retry_count: 0,
          }),
        );
        return;
      }

      if (method === "GET" && path === "/organizations/org_test/environments") {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify([
            {
              id: "env_123",
              org_id: "org_test",
              state: "running",
              desired_state: "running",
              type: "sandbox",
              created_by: "user_1",
              template_id: "tmpl_1",
              recoverable: true,
              state_changed_at: "2026-01-01T00:00:00Z",
              cpu_millicores: 2000,
              memory_mb: 4096,
              disk_gb: 20,
              total_running_seconds: 0,
              created_at: "2026-01-01T00:00:00Z",
              retry_count: 0,
            },
          ]),
        );
        return;
      }

      if (
        method === "POST" &&
        path === "/organizations/org_test/environments/env_123/exec"
      ) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({ exit_code: 0, stdout: "hello\n", stderr: "" }),
        );
        return;
      }

      if (
        method === "POST" &&
        path === "/organizations/org_test/environments/env_123/stop"
      ) {
        res.writeHead(204);
        res.end();
        return;
      }

      if (
        method === "DELETE" &&
        path === "/organizations/org_test/environments/env_123"
      ) {
        res.writeHead(204);
        res.end();
        return;
      }

      // Templates
      if (method === "POST" && path === "/organizations/org_test/templates") {
        res.writeHead(201, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            id: "tmpl_1",
            org_id: "org_test",
            name: "Go",
            type: "sandbox",
            build_status: "ready",
            default_cpu_millicores: 2000,
            default_memory_mb: 4096,
            default_disk_gb: 20,
            official: false,
            created_at: "2026-01-01T00:00:00Z",
          }),
        );
        return;
      }

      if (
        method === "GET" &&
        path === "/organizations/org_test/templates/tmpl_1"
      ) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            id: "tmpl_1",
            org_id: "org_test",
            name: "Go",
            type: "sandbox",
            build_status: "ready",
            default_cpu_millicores: 2000,
            default_memory_mb: 4096,
            default_disk_gb: 20,
            official: false,
            created_at: "2026-01-01T00:00:00Z",
          }),
        );
        return;
      }

      // API Keys
      if (method === "POST" && path === "/organizations/org_test/api-keys") {
        res.writeHead(201, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            id: "key_1",
            user_id: "user_1",
            org_id: "org_test",
            name: "test",
            key: "cw_abc123",
            key_prefix: "cw_abc",
            scopes: [],
            created_at: "2026-01-01T00:00:00Z",
          }),
        );
        return;
      }

      if (method === "GET" && path === "/organizations/org_test/api-keys") {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify([
            {
              id: "key_1",
              user_id: "user_1",
              org_id: "org_test",
              name: "test",
              key_prefix: "cw_abc",
              scopes: [],
              created_at: "2026-01-01T00:00:00Z",
            },
          ]),
        );
        return;
      }

      // Secrets
      if (method === "GET" && path === "/secrets") {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ secrets: [{ key: "DB_URL" }] }));
        return;
      }

      if (method === "PUT" && path === "/secrets") {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ status: "ok" }));
        return;
      }

      if (method === "DELETE" && path.startsWith("/secrets/")) {
        res.writeHead(204);
        res.end();
        return;
      }

      // 404
      res.writeHead(404);
      res.end(JSON.stringify({ detail: "not found" }));
    },
  );

  await new Promise<void>((resolve) => server.listen(0, resolve));
  const addr = server.address() as { port: number };
  const baseUrl = `http://127.0.0.1:${addr.port}`;

  try {
    const cw = new Codewire({
      apiKey: "test-key",
      orgId: "org_test",
      baseUrl,
    });

    // Test: Create environment returns Environment wrapper
    console.log("Test: Environment CRUD");
    const env = await cw.environments.create({ template_slug: "python" });
    assert(env instanceof Environment, "create returns Environment instance");
    assert(env.id === "env_123", "env.id matches");
    assert(env.state === "running", "env.state matches");

    // Test: Get environment
    const env2 = await cw.environments.get("env_123");
    assert(env2.id === "env_123", "get returns correct env");

    // Test: List environments
    const envs = await cw.environments.list();
    assert(envs.length === 1, "list returns 1 env");
    assert(envs[0] instanceof Environment, "list items are Environment instances");

    // Test: Exec via wrapper
    console.log("Test: Environment exec");
    const result = await env.exec({ command: ["echo", "hello"] });
    assert(result.exit_code === 0, "exec exit code 0");
    assert(result.stdout === "hello\n", "exec stdout matches");

    // Test: Stop via wrapper
    console.log("Test: Environment stop/remove");
    await env.stop();
    await env.remove();

    // Test: Templates
    console.log("Test: Templates");
    const tmpl = await cw.templates.create({
      name: "Go",
      type: "sandbox",
    });
    assert(tmpl.id === "tmpl_1", "template created");

    const tmpl2 = await cw.templates.get("tmpl_1");
    assert(tmpl2.name === "Go", "template get works");

    // Test: API Keys
    console.log("Test: API Keys");
    const key = await cw.apiKeys.create({ name: "test" });
    assert(key.key === "cw_abc123", "api key created");

    const keys = await cw.apiKeys.list();
    assert(keys.length === 1, "api keys listed");

    // Test: Secrets
    console.log("Test: Secrets");
    const secrets = await cw.secrets.list();
    assert(secrets.length === 1, "secrets listed");
    assert(secrets[0].key === "DB_URL", "secret key matches");

    await cw.secrets.set("DB_URL", "postgres://...");
    await cw.secrets.delete("DB_URL");

    // Test: Error handling
    console.log("Test: Error handling");
    try {
      const cw2 = new Codewire({
        apiKey: "wrong-key",
        orgId: "org_test",
        baseUrl,
      });
      await cw2.environments.list();
      assert(false, "should have thrown");
    } catch (e) {
      assert(e instanceof CodewireError, "throws CodewireError");
      assert((e as CodewireError).status === 401, "status is 401");
    }

    // Test: Missing orgId
    console.log("Test: Missing orgId");
    try {
      const cw3 = new Codewire({ apiKey: "test-key", baseUrl });
      await cw3.environments.list();
      assert(false, "should have thrown for missing orgId");
    } catch (e) {
      assert(
        (e as Error).message.includes("orgId is required"),
        "error mentions orgId",
      );
    }
  } finally {
    server.close();
  }

  console.log(`\n${passed} passed, ${failed} failed`);
  if (failed > 0) process.exit(1);
}

runTests().catch((e) => {
  console.error(e);
  process.exit(1);
});
