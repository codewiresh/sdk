from __future__ import annotations

from typing import Any

from codewire.http import HttpClient


class Environments:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, org_id: str, **kwargs: Any) -> dict:
        return self._http.post(f"/organizations/{org_id}/environments", kwargs)

    def list(
        self,
        org_id: str,
        *,
        type: str | None = None,
        state: str | None = None,
    ) -> list[dict]:
        params = {}
        if type:
            params["type"] = type
        if state:
            params["state"] = state
        qs = "&".join(f"{k}={v}" for k, v in params.items())
        path = f"/organizations/{org_id}/environments"
        if qs:
            path += f"?{qs}"
        return self._http.get(path)

    def get(self, org_id: str, env_id: str) -> dict:
        return self._http.get(f"/organizations/{org_id}/environments/{env_id}")

    def delete(self, org_id: str, env_id: str) -> None:
        self._http.delete(f"/organizations/{org_id}/environments/{env_id}")

    def stop(self, org_id: str, env_id: str) -> dict:
        return self._http.post(
            f"/organizations/{org_id}/environments/{env_id}/stop"
        )

    def start(self, org_id: str, env_id: str) -> dict:
        return self._http.post(
            f"/organizations/{org_id}/environments/{env_id}/start"
        )


class Templates:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, org_id: str, **kwargs: Any) -> dict:
        return self._http.post(f"/organizations/{org_id}/templates", kwargs)

    def list(self, org_id: str, *, type: str | None = None) -> list[dict]:
        path = f"/organizations/{org_id}/templates"
        if type:
            path += f"?type={type}"
        return self._http.get(path)

    def delete(self, org_id: str, template_id: str) -> None:
        self._http.delete(f"/organizations/{org_id}/templates/{template_id}")


class APIKeys:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, org_id: str, **kwargs: Any) -> dict:
        return self._http.post(f"/organizations/{org_id}/api-keys", kwargs)

    def list(self, org_id: str) -> list[dict]:
        return self._http.get(f"/organizations/{org_id}/api-keys")

    def delete(self, org_id: str, key_id: str) -> None:
        self._http.delete(f"/organizations/{org_id}/api-keys/{key_id}")
