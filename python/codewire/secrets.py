from __future__ import annotations

from typing import Any
from urllib.parse import quote

from codewire.http import HttpClient


class Secrets:
    def __init__(self, http: HttpClient):
        self._http = http

    def list(self) -> list[dict]:
        org_id = self._http.resolve_org_id()
        resp = self._http.get(f"/secrets?org_id={quote(org_id)}")
        return resp.get("secrets", []) if isinstance(resp, dict) else resp

    def set(self, key: str, value: str) -> None:
        org_id = self._http.resolve_org_id()
        self._http.put("/secrets", {"org_id": org_id, "key": key, "value": value})

    def delete(self, key: str) -> None:
        org_id = self._http.resolve_org_id()
        self._http.delete(f"/secrets/{quote(key)}?org_id={quote(org_id)}")

    def list_user(self) -> list[dict]:
        resp = self._http.get("/user/secrets")
        return resp.get("secrets", []) if isinstance(resp, dict) else resp

    def set_user(self, key: str, value: str) -> None:
        self._http.put("/user/secrets", {"key": key, "value": value})

    def delete_user(self, key: str) -> None:
        self._http.delete(f"/user/secrets/{quote(key)}")


class SecretProjects:
    def __init__(self, http: HttpClient):
        self._http = http

    def create(self, name: str) -> dict:
        return self._http.post(
            self._http.org_path("/secret-projects"), {"name": name}
        )

    def list(self) -> list[dict]:
        return self._http.get(self._http.org_path("/secret-projects"))

    def delete(self, project_id: str) -> None:
        self._http.delete(
            self._http.org_path(f"/secret-projects/{project_id}")
        )

    def list_secrets(self, project_id: str) -> list[dict]:
        resp = self._http.get(
            self._http.org_path(f"/secret-projects/{project_id}/secrets")
        )
        return resp.get("secrets", []) if isinstance(resp, dict) else resp

    def set_secret(self, project_id: str, key: str, value: str) -> None:
        self._http.put(
            self._http.org_path(f"/secret-projects/{project_id}/secrets"),
            {"key": key, "value": value},
        )

    def delete_secret(self, project_id: str, key: str) -> None:
        self._http.delete(
            self._http.org_path(
                f"/secret-projects/{project_id}/secrets/{quote(key)}"
            )
        )
