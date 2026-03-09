#!/usr/bin/env python3
"""Filter the full Codewire OpenAPI spec down to SDK-relevant endpoints."""

import json
import sys
from collections import deque


TAG_WHITELIST = {
    "environments",
    "templates",
    "api-keys",
    "secrets",
    "secret-projects",
    "ports",
}


def collect_refs(obj, refs=None):
    """Recursively collect all $ref targets from a JSON structure."""
    if refs is None:
        refs = set()
    if isinstance(obj, dict):
        if "$ref" in obj:
            ref = obj["$ref"]
            if ref.startswith("#/components/schemas/"):
                refs.add(ref.split("/")[-1])
        for v in obj.values():
            collect_refs(v, refs)
    elif isinstance(obj, list):
        for item in obj:
            collect_refs(item, refs)
    return refs


def strip_schema_property(obj):
    """Remove Huma's $schema properties from schema objects."""
    if isinstance(obj, dict):
        obj.pop("$schema", None)
        props = obj.get("properties", {})
        props.pop("$schema", None)
        req = obj.get("required")
        if isinstance(req, list) and "$schema" in req:
            req.remove("$schema")
        for v in obj.values():
            strip_schema_property(v)
    elif isinstance(obj, list):
        for item in obj:
            strip_schema_property(item)


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input-openapi> <output-json>", file=sys.stderr)
        sys.exit(1)

    src, dst = sys.argv[1], sys.argv[2]

    with open(src) as f:
        spec = json.load(f)

    # Filter paths by tag whitelist
    filtered_paths = {}
    for path, methods in spec.get("paths", {}).items():
        filtered_methods = {}
        for method, op in methods.items():
            if isinstance(op, dict) and set(op.get("tags", [])) & TAG_WHITELIST:
                filtered_methods[method] = op
        if filtered_methods:
            filtered_paths[path] = filtered_methods

    # Collect all transitively referenced schemas
    direct_refs = collect_refs(filtered_paths)
    all_schemas = spec.get("components", {}).get("schemas", {})

    resolved = set()
    queue = deque(direct_refs)
    while queue:
        name = queue.popleft()
        if name in resolved:
            continue
        resolved.add(name)
        schema = all_schemas.get(name, {})
        for dep in collect_refs(schema):
            if dep not in resolved:
                queue.append(dep)

    # Build filtered spec
    filtered_schemas = {k: v for k, v in all_schemas.items() if k in resolved}

    out = {
        "openapi": spec.get("openapi", "3.0.0"),
        "info": {
            "title": "Codewire SDK API",
            "version": spec.get("info", {}).get("version", "1.0.0"),
        },
        "paths": filtered_paths,
        "components": {
            "schemas": filtered_schemas,
            "securitySchemes": {
                "BearerAuth": {
                    "type": "http",
                    "scheme": "bearer",
                }
            },
        },
    }

    strip_schema_property(out)

    with open(dst, "w") as f:
        json.dump(out, f, indent=2, sort_keys=False)

    print(
        f"Filtered: {len(filtered_paths)} paths, {len(filtered_schemas)} schemas",
        file=sys.stderr,
    )


if __name__ == "__main__":
    main()
