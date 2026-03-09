from __future__ import annotations

from typing import Any
from urllib.parse import quote

import httpx

from codewire.errors import CodewireError


class HttpClient:
    def __init__(self, base_url: str, api_key: str, org_id: str):
        self._client = httpx.Client(
            base_url=base_url,
            headers={
                "Authorization": f"Bearer {api_key}",
                "Content-Type": "application/json",
            },
        )
        self.org_id = org_id

    def resolve_org_id(self) -> str:
        if not self.org_id:
            raise ValueError(
                "org_id is required — pass org_id to Codewire() "
                "or set CODEWIRE_ORG_ID"
            )
        return self.org_id

    def org_path(self, pattern: str) -> str:
        return f"/organizations/{self.resolve_org_id()}{pattern}"

    def request(self, method: str, path: str, body: Any = None) -> Any:
        resp = self._client.request(method, path, json=body)
        if resp.status_code >= 400:
            detail = resp.reason_phrase or ""
            try:
                data = resp.json()
                detail = data.get("detail") or data.get("message") or detail
            except Exception:
                pass
            raise CodewireError(resp.status_code, detail)
        if resp.status_code == 204:
            return None
        return resp.json()

    def request_raw(
        self,
        method: str,
        path: str,
        data: bytes | None = None,
        content_type: str | None = None,
    ) -> httpx.Response:
        headers: dict[str, str] = {}
        if content_type:
            headers["Content-Type"] = content_type
        resp = self._client.request(method, path, content=data, headers=headers)
        if resp.status_code >= 400:
            detail = resp.reason_phrase or ""
            try:
                d = resp.json()
                detail = d.get("detail") or d.get("message") or detail
            except Exception:
                pass
            raise CodewireError(resp.status_code, detail)
        return resp

    def get(self, path: str) -> Any:
        return self.request("GET", path)

    def post(self, path: str, body: Any = None) -> Any:
        return self.request("POST", path, body)

    def put(self, path: str, body: Any = None) -> Any:
        return self.request("PUT", path, body)

    def patch(self, path: str, body: Any = None) -> Any:
        return self.request("PATCH", path, body)

    def delete(self, path: str) -> Any:
        return self.request("DELETE", path)

    def close(self) -> None:
        self._client.close()
