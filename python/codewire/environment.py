from __future__ import annotations

from typing import Any, TYPE_CHECKING
from urllib.parse import quote

if TYPE_CHECKING:
    from codewire.http import HttpClient


class Environment:
    """Wrapper around environment API data with convenience methods."""

    def __init__(self, data: dict, http: HttpClient):
        self._http = http
        self._org_id = http.org_id
        self._data = data
        # Populate public attributes from API response
        for key, value in data.items():
            setattr(self, key, value)

    @property
    def _path_prefix(self) -> str:
        return f"/organizations/{self._org_id}/environments/{self.id}"

    def exec(self, **kwargs: Any) -> dict:
        return self._http.post(f"{self._path_prefix}/exec", kwargs)

    def start(self) -> None:
        self._http.post(f"{self._path_prefix}/start")

    def stop(self) -> None:
        self._http.post(f"{self._path_prefix}/stop")

    def remove(self) -> None:
        self._http.delete(self._path_prefix)

    def upload(self, src: bytes, remote_path: str) -> None:
        qs = f"?path={quote(remote_path)}"
        self._http.request_raw(
            "POST",
            f"{self._path_prefix}/files/upload{qs}",
            data=src,
            content_type="application/octet-stream",
        )

    def download(self, remote_path: str) -> bytes:
        qs = f"?path={quote(remote_path)}"
        resp = self._http.request_raw(
            "GET",
            f"{self._path_prefix}/files/download{qs}",
        )
        return resp.content

    def list_files(self, path: str = "/") -> list[dict]:
        qs = f"?path={quote(path)}"
        return self._http.get(f"{self._path_prefix}/files{qs}")

    def list_ports(self) -> list[dict]:
        return self._http.get(f"{self._path_prefix}/ports")

    def create_port(self, **kwargs: Any) -> dict:
        return self._http.post(f"{self._path_prefix}/ports", kwargs)

    def delete_port(self, port_id: str) -> None:
        self._http.delete(f"{self._path_prefix}/ports/{port_id}")

    def __repr__(self) -> str:
        state = getattr(self, "state", "?")
        env_id = getattr(self, "id", "?")
        return f"<Environment id={env_id} state={state}>"
