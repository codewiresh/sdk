from __future__ import annotations

from typing import Any

import httpx

from codewire.errors import CodewireError


class HttpClient:
    def __init__(self, base_url: str, api_key: str):
        self._client = httpx.Client(
            base_url=base_url,
            headers={
                "Authorization": f"Bearer {api_key}",
                "Content-Type": "application/json",
            },
        )

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

    def get(self, path: str) -> Any:
        return self.request("GET", path)

    def post(self, path: str, body: Any = None) -> Any:
        return self.request("POST", path, body)

    def delete(self, path: str) -> Any:
        return self.request("DELETE", path)

    def close(self) -> None:
        self._client.close()
