from __future__ import annotations

from typing import Any
from urllib.parse import quote

from codewire.http import HttpClient
from codewire.environment import Environment
from codewire.models import (
    CreateAPIKeyBody,
    CreateAPIKeyResponse,
    CreateEnvironmentBody,
    CreatePortBody,
    CreateTemplateBody,
    UpdateTemplateBody,
    Environment as EnvironmentModel,
    EnvironmentTemplate,
    EnvironmentPort,
    ExecBody,
    ExecResult,
    FileEntry,
    APIKey,
)


class Environments:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, **kwargs: Any) -> Environment:
        data = self._http.post(self._http.org_path("/environments"), kwargs)
        return Environment(data, self._http)

    def list(
        self,
        *,
        type: str | None = None,
        state: str | None = None,
    ) -> list[Environment]:
        params: dict[str, str] = {}
        if type:
            params["type"] = type
        if state:
            params["state"] = state
        qs = "&".join(f"{k}={v}" for k, v in params.items())
        path = self._http.org_path("/environments")
        if qs:
            path += f"?{qs}"
        items = self._http.get(path)
        return [Environment(d, self._http) for d in items]

    def get(self, env_id: str) -> Environment:
        data = self._http.get(
            self._http.org_path(f"/environments/{env_id}")
        )
        return Environment(data, self._http)

    def delete(self, env_id: str) -> None:
        self._http.delete(self._http.org_path(f"/environments/{env_id}"))

    def stop(self, env_id: str) -> None:
        self._http.post(self._http.org_path(f"/environments/{env_id}/stop"))

    def start(self, env_id: str) -> None:
        self._http.post(self._http.org_path(f"/environments/{env_id}/start"))

    def exec(self, env_id: str, **kwargs: Any) -> dict:
        return self._http.post(
            self._http.org_path(f"/environments/{env_id}/exec"), kwargs
        )

    def list_files(self, env_id: str, path: str | None = None) -> list[dict]:
        qs = f"?path={quote(path)}" if path else ""
        return self._http.get(
            self._http.org_path(f"/environments/{env_id}/files{qs}")
        )

    def upload_file(
        self, env_id: str, src: bytes, remote_path: str
    ) -> None:
        qs = f"?path={quote(remote_path)}"
        self._http.request_raw(
            "POST",
            self._http.org_path(f"/environments/{env_id}/files/upload{qs}"),
            data=src,
            content_type="application/octet-stream",
        )

    def download_file(self, env_id: str, remote_path: str) -> bytes:
        qs = f"?path={quote(remote_path)}"
        resp = self._http.request_raw(
            "GET",
            self._http.org_path(f"/environments/{env_id}/files/download{qs}"),
        )
        return resp.content

    def list_ports(self, env_id: str) -> list[dict]:
        return self._http.get(
            self._http.org_path(f"/environments/{env_id}/ports")
        )

    def create_port(self, env_id: str, **kwargs: Any) -> dict:
        return self._http.post(
            self._http.org_path(f"/environments/{env_id}/ports"), kwargs
        )

    def delete_port(self, env_id: str, port_id: str) -> None:
        self._http.delete(
            self._http.org_path(f"/environments/{env_id}/ports/{port_id}")
        )


class Templates:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, **kwargs: Any) -> dict:
        return self._http.post(self._http.org_path("/templates"), kwargs)

    def list(self, *, type: str | None = None) -> list[dict]:
        path = self._http.org_path("/templates")
        if type:
            path += f"?type={type}"
        return self._http.get(path)

    def get(self, template_id: str) -> dict:
        return self._http.get(
            self._http.org_path(f"/templates/{template_id}")
        )

    def update(self, template_id: str, **kwargs: Any) -> dict:
        return self._http.patch(
            self._http.org_path(f"/templates/{template_id}"), kwargs
        )

    def delete(self, template_id: str) -> None:
        self._http.delete(self._http.org_path(f"/templates/{template_id}"))


class APIKeys:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, **kwargs: Any) -> dict:
        return self._http.post(self._http.org_path("/api-keys"), kwargs)

    def list(self) -> list[dict]:
        return self._http.get(self._http.org_path("/api-keys"))

    def delete(self, key_id: str) -> None:
        self._http.delete(self._http.org_path(f"/api-keys/{key_id}"))
