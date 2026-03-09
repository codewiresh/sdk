"""Tests for the Codewire Python SDK using a mock HTTP server."""

import json
import threading
import unittest
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

from codewire import Codewire, CodewireError
from codewire.environment import Environment


ENV_JSON = {
    "id": "env_123",
    "org_id": "org_test",
    "state": "running",
    "desired_state": "running",
    "type": "sandbox",
    "created_by": "user_1",
    "template_id": "tmpl_1",
    "recoverable": True,
    "state_changed_at": "2026-01-01T00:00:00Z",
    "cpu_millicores": 2000,
    "memory_mb": 4096,
    "disk_gb": 20,
    "total_running_seconds": 0,
    "created_at": "2026-01-01T00:00:00Z",
    "retry_count": 0,
}

TMPL_JSON = {
    "id": "tmpl_1",
    "org_id": "org_test",
    "name": "Go",
    "type": "sandbox",
    "build_status": "ready",
    "default_cpu_millicores": 2000,
    "default_memory_mb": 4096,
    "default_disk_gb": 20,
    "official": False,
    "created_at": "2026-01-01T00:00:00Z",
}


class MockHandler(BaseHTTPRequestHandler):
    def log_message(self, *args):
        pass  # suppress output

    def _respond(self, code, body=None):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        if body is not None:
            self.wfile.write(json.dumps(body).encode())

    def _check_auth(self):
        auth = self.headers.get("Authorization")
        if auth != "Bearer test-key":
            self._respond(401, {"detail": "unauthorized"})
            return False
        return True

    def do_GET(self):
        if not self._check_auth():
            return
        parsed = urlparse(self.path)
        path = parsed.path

        if path == "/organizations/org_test/environments":
            self._respond(200, [ENV_JSON])
        elif path == "/organizations/org_test/environments/env_123":
            self._respond(200, ENV_JSON)
        elif path.endswith("/files"):
            self._respond(200, [{"name": "test.txt", "is_dir": False, "size": 10, "mode": "-rw-r--r--", "mod_time": "2026-01-01T00:00:00Z"}])
        elif path.endswith("/ports"):
            self._respond(200, [{"id": "port_1", "environment_id": "env_123", "port": 8080, "label": "web", "access": "public", "created_at": "2026-01-01T00:00:00Z"}])
        elif path == "/organizations/org_test/templates":
            self._respond(200, [TMPL_JSON])
        elif path == "/organizations/org_test/templates/tmpl_1":
            self._respond(200, TMPL_JSON)
        elif path == "/organizations/org_test/api-keys":
            self._respond(200, [{"id": "key_1", "user_id": "user_1", "org_id": "org_test", "name": "test", "key_prefix": "cw_abc", "scopes": [], "created_at": "2026-01-01T00:00:00Z"}])
        elif path == "/secrets":
            self._respond(200, {"secrets": [{"key": "DB_URL"}]})
        elif path == "/user/secrets":
            self._respond(200, {"secrets": [{"key": "TOKEN"}]})
        elif path == "/organizations/org_test/secret-projects":
            self._respond(200, [{"id": "proj_1", "org_id": "org_test", "name": "myapp", "created_at": "2026-01-01T00:00:00Z"}])
        else:
            self._respond(404, {"detail": "not found"})

    def do_POST(self):
        if not self._check_auth():
            return
        path = urlparse(self.path).path

        if path == "/organizations/org_test/environments":
            self._respond(201, ENV_JSON)
        elif path.endswith("/exec"):
            self._respond(200, {"exit_code": 0, "stdout": "hello\n", "stderr": ""})
        elif path.endswith("/stop") or path.endswith("/start"):
            self._respond(204)
        elif path.endswith("/ports"):
            self._respond(201, {"id": "port_1", "environment_id": "env_123", "port": 8080, "label": "web", "access": "public", "created_at": "2026-01-01T00:00:00Z"})
        elif path == "/organizations/org_test/templates":
            self._respond(201, TMPL_JSON)
        elif path == "/organizations/org_test/api-keys":
            self._respond(201, {"id": "key_1", "user_id": "user_1", "org_id": "org_test", "name": "test", "key": "cw_abc123", "key_prefix": "cw_abc", "scopes": [], "created_at": "2026-01-01T00:00:00Z"})
        elif path == "/organizations/org_test/secret-projects":
            self._respond(201, {"id": "proj_1", "org_id": "org_test", "name": "myapp", "created_at": "2026-01-01T00:00:00Z"})
        else:
            self._respond(404, {"detail": "not found"})

    def do_PUT(self):
        if not self._check_auth():
            return
        self._respond(200, {"status": "ok"})

    def do_PATCH(self):
        if not self._check_auth():
            return
        self._respond(200, TMPL_JSON)

    def do_DELETE(self):
        if not self._check_auth():
            return
        self._respond(204)


class SDKTestCase(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server = HTTPServer(("127.0.0.1", 0), MockHandler)
        cls.port = cls.server.server_address[1]
        cls.thread = threading.Thread(target=cls.server.serve_forever)
        cls.thread.daemon = True
        cls.thread.start()
        cls.base_url = f"http://127.0.0.1:{cls.port}"
        cls.cw = Codewire(
            api_key="test-key",
            org_id="org_test",
            base_url=cls.base_url,
        )

    @classmethod
    def tearDownClass(cls):
        cls.server.shutdown()
        cls.cw.close()

    def test_environment_create(self):
        env = self.cw.environments.create(template_slug="python")
        self.assertIsInstance(env, Environment)
        self.assertEqual(env.id, "env_123")
        self.assertEqual(env.state, "running")

    def test_environment_get(self):
        env = self.cw.environments.get("env_123")
        self.assertEqual(env.id, "env_123")

    def test_environment_list(self):
        envs = self.cw.environments.list()
        self.assertEqual(len(envs), 1)
        self.assertIsInstance(envs[0], Environment)

    def test_environment_exec(self):
        env = self.cw.environments.create(template_slug="python")
        result = env.exec(command=["echo", "hello"])
        self.assertEqual(result["exit_code"], 0)
        self.assertEqual(result["stdout"], "hello\n")

    def test_environment_stop(self):
        env = self.cw.environments.create(template_slug="python")
        env.stop()  # should not raise

    def test_environment_remove(self):
        env = self.cw.environments.create(template_slug="python")
        env.remove()  # should not raise

    def test_environment_list_files(self):
        env = self.cw.environments.create(template_slug="python")
        files = env.list_files("/workspace")
        self.assertEqual(len(files), 1)
        self.assertEqual(files[0]["name"], "test.txt")

    def test_environment_list_ports(self):
        env = self.cw.environments.create(template_slug="python")
        ports = env.list_ports()
        self.assertEqual(len(ports), 1)
        self.assertEqual(ports[0]["port"], 8080)

    def test_templates_crud(self):
        tmpl = self.cw.templates.create(name="Go", type="sandbox")
        self.assertEqual(tmpl["name"], "Go")

        tmpl2 = self.cw.templates.get("tmpl_1")
        self.assertEqual(tmpl2["id"], "tmpl_1")

        tmpls = self.cw.templates.list()
        self.assertEqual(len(tmpls), 1)

        updated = self.cw.templates.update("tmpl_1", name="Golang")
        self.assertEqual(updated["name"], "Go")  # mock returns same

        self.cw.templates.delete("tmpl_1")

    def test_api_keys(self):
        key = self.cw.api_keys.create(name="test")
        self.assertEqual(key["key"], "cw_abc123")

        keys = self.cw.api_keys.list()
        self.assertEqual(len(keys), 1)

        self.cw.api_keys.delete("key_1")

    def test_secrets(self):
        secrets = self.cw.secrets.list()
        self.assertEqual(len(secrets), 1)
        self.assertEqual(secrets[0]["key"], "DB_URL")

        self.cw.secrets.set("DB_URL", "postgres://...")
        self.cw.secrets.delete("DB_URL")

    def test_user_secrets(self):
        secrets = self.cw.secrets.list_user()
        self.assertEqual(len(secrets), 1)

        self.cw.secrets.set_user("TOKEN", "abc")
        self.cw.secrets.delete_user("TOKEN")

    def test_secret_projects(self):
        proj = self.cw.secret_projects.create("myapp")
        self.assertEqual(proj["name"], "myapp")

        projects = self.cw.secret_projects.list()
        self.assertEqual(len(projects), 1)

        self.cw.secret_projects.delete("proj_1")

    def test_error_handling(self):
        bad_cw = Codewire(
            api_key="wrong-key",
            org_id="org_test",
            base_url=self.base_url,
        )
        with self.assertRaises(CodewireError) as cm:
            bad_cw.environments.list()
        self.assertEqual(cm.exception.status, 401)
        bad_cw.close()

    def test_missing_org_id(self):
        no_org = Codewire(api_key="test-key", base_url=self.base_url)
        with self.assertRaises(ValueError) as cm:
            no_org.environments.list()
        self.assertIn("org_id is required", str(cm.exception))
        no_org.close()

    def test_environment_repr(self):
        env = self.cw.environments.create(template_slug="python")
        r = repr(env)
        self.assertIn("env_123", r)
        self.assertIn("running", r)


if __name__ == "__main__":
    unittest.main()
