from __future__ import annotations

import os

from codewire.http import HttpClient
from codewire.resources import APIKeys, Environments, Templates
from codewire.secrets import Secrets, SecretProjects

DEFAULT_BASE_URL = "https://api.codewire.sh"


class Codewire:
    def __init__(
        self,
        api_key: str | None = None,
        base_url: str | None = None,
        org_id: str | None = None,
    ):
        key = api_key or os.environ.get("CODEWIRE_API_KEY", "")
        if not key:
            raise ValueError(
                "Codewire API key is required. "
                "Pass api_key or set CODEWIRE_API_KEY."
            )
        url = (base_url or DEFAULT_BASE_URL).rstrip("/")
        resolved_org = org_id or os.environ.get("CODEWIRE_ORG_ID", "")
        self._http = HttpClient(url, key, resolved_org)
        self.environments = Environments(self._http)
        self.templates = Templates(self._http)
        self.api_keys = APIKeys(self._http)
        self.secrets = Secrets(self._http)
        self.secret_projects = SecretProjects(self._http)

    def close(self) -> None:
        self._http.close()

    def __enter__(self) -> Codewire:
        return self

    def __exit__(self, *args: object) -> None:
        self.close()
