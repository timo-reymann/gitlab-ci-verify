#!/usr/bin/env bash
#
# Generate NOTICE file from syft SBOM output.
# Uses SPDX license identifiers and fetches canonical license texts from spdx.org.
#
# Requirements: syft, python3
#
# Usage: ./generate-notice.sh [syft-args...] > NOTICE
#   e.g.: ./generate-notice.sh dir:.
#         ./generate-notice.sh alpine:latest

set -euo pipefail

if [ $# -eq 0 ]; then
    echo "Usage: $0 <syft-target> [syft-args...]" >&2
    echo "  e.g.: $0 dir:." >&2
    exit 1
fi

SBOM_JSON=$(syft "$@" -o syft-json 2>/dev/null)

echo "$SBOM_JSON" | python3 -c '
import json
import sys
import urllib.request
import os

SPDX_LICENSE_TEXT_URL = "https://raw.githubusercontent.com/spdx/license-list-data/main/text/{}.txt"

# Map common non-SPDX license names to their SPDX equivalents
NON_SPDX_ALIASES = {
    "Apache 2.0": "Apache-2.0",
    "Apache License, Version 2.0": "Apache-2.0",
    "Apache Software License": "Apache-2.0",
    "MIT License": "MIT",
    "MIT/X11": "MIT",
    "BSD License": "BSD-3-Clause",
    "Dual License": None,  # ambiguous, skip
    "UNKNOWN": None,
    "ZPL 2.1": "ZPL-2.1",
    "GooglePatentsFile": None,
}

def is_sha256(value):
    """Check if a value is a sha256 hash (not a license)."""
    return value.startswith("sha256:")

def normalize_license(value):
    """Normalize a license value to an SPDX ID, or None to skip."""
    if is_sha256(value):
        return None
    if value in NON_SPDX_ALIASES:
        return NON_SPDX_ALIASES[value]
    return value

def fetch_spdx_license_text(license_id):
    """Fetch canonical license text from spdx.org license-list-data."""
    url = SPDX_LICENSE_TEXT_URL.format(license_id)
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            return resp.read().decode("utf-8").strip()
    except Exception:
        return None

def main():
    data = json.load(sys.stdin)
    artifacts = data.get("artifacts", [])

    # Collect unique packages with their licenses
    # Use (name, version) as key to deduplicate
    packages = {}
    for a in artifacts:
        key = (a["name"], a.get("version", "unknown"))
        if key in packages:
            continue

        spdx_ids = set()
        for lic in a.get("licenses", []):
            expr = lic.get("spdxExpression", "")
            if expr:
                for part in expr.replace("(", "").replace(")", "").split():
                    part = part.strip()
                    if part and part not in ("AND", "OR", "WITH"):
                        normalized = normalize_license(part)
                        if normalized:
                            spdx_ids.add(normalized)
            elif lic.get("value"):
                normalized = normalize_license(lic["value"])
                if normalized:
                    spdx_ids.add(normalized)

        if spdx_ids:
            packages[key] = sorted(spdx_ids)

    if not packages:
        print("No packages with license information found.", file=sys.stderr)
        sys.exit(1)

    # Collect all unique SPDX IDs and fetch their texts
    all_license_ids = set()
    for spdx_ids in packages.values():
        all_license_ids.update(spdx_ids)

    license_texts = {}
    cache_dir = os.path.join(os.environ.get("TMPDIR", "/tmp"), "spdx-license-cache")
    os.makedirs(cache_dir, exist_ok=True)

    for lid in sorted(all_license_ids):
        cache_file = os.path.join(cache_dir, f"{lid}.txt")
        if os.path.exists(cache_file):
            with open(cache_file, "r") as f:
                license_texts[lid] = f.read()
            continue

        text = fetch_spdx_license_text(lid)
        if text:
            license_texts[lid] = text
            with open(cache_file, "w") as f:
                f.write(text)
        else:
            print(f"Warning: Could not fetch license text for {lid}", file=sys.stderr)

    # Generate NOTICE
    print("This software includes external packages and source code.")
    print("The applicable license information is listed below:")

    for (name, version), spdx_ids in sorted(packages.items()):
        print()
        print("----")
        print()
        print(f"Package: {name}:{version}")
        for lid in spdx_ids:
            print(f"License: {lid}")

        for lid in spdx_ids:
            text = license_texts.get(lid)
            if text:
                print()
                print(text)

main()
'
