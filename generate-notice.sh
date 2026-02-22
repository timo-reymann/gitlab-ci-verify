#!/usr/bin/env bash
#
# Generate NOTICE file from syft SBOM output.
# Reads actual license files from package sources on disk when available,
# falls back to canonical SPDX license texts from spdx.org.
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

echo "$SBOM_JSON" | SCAN_DIR="$(pwd)" python3 -c '
import json
import sys
import subprocess
import urllib.request
import os
import glob as globmod

SPDX_LICENSE_TEXT_URL = "https://raw.githubusercontent.com/spdx/license-list-data/main/text/{}.txt"

# Map common non-SPDX license names to their SPDX equivalents
NON_SPDX_ALIASES = {
    "Apache 2.0": "Apache-2.0",
    "Apache License, Version 2.0": "Apache-2.0",
    "Apache Software License": "Apache-2.0",
    "MIT License": "MIT",
    "MIT/X11": "MIT",
    "BSD License": "BSD-3-Clause",
    "Dual License": None,
    "UNKNOWN": None,
    "ZPL 2.1": "ZPL-2.1",
    "GooglePatentsFile": None,
}

LICENSE_FILE_PATTERNS = ["LICENSE*", "LICENCE*", "COPYING*", "COPYRIGHT*"]

def is_sha256(value):
    return value.startswith("sha256:")

def normalize_license(value):
    if is_sha256(value):
        return None
    if value in NON_SPDX_ALIASES:
        return NON_SPDX_ALIASES[value]
    return value

def go_module_cache_encode(name):
    """Encode a Go module name for the module cache (uppercase -> !lowercase)."""
    result = []
    for ch in name:
        if ch.isupper():
            result.append("!" + ch.lower())
        else:
            result.append(ch)
    return "".join(result)

def get_gomodcache():
    """Get the Go module cache directory."""
    try:
        result = subprocess.run(
            ["go", "env", "GOMODCACHE"],
            capture_output=True, text=True, timeout=5
        )
        path = result.stdout.strip()
        if path and os.path.isdir(path):
            return path
    except Exception:
        pass
    return None

def find_license_file_in_dir(directory):
    """Find a license file in a directory, return its content or None."""
    if not os.path.isdir(directory):
        return None
    for pattern in LICENSE_FILE_PATTERNS:
        matches = globmod.glob(os.path.join(directory, pattern))
        for match in sorted(matches):
            if os.path.isfile(match):
                try:
                    with open(match, "r", errors="replace") as f:
                        return f.read().strip()
                except Exception:
                    continue
    return None

def find_license_from_syft_locations(locations, scan_dir):
    """Try to read a license file from syft-reported locations."""
    for loc in locations:
        path = loc.get("path", "") or loc.get("accessPath", "")
        if not path:
            continue
        # Paths from syft are relative to the scan root
        full_path = os.path.join(scan_dir, path.lstrip("/"))
        if os.path.isfile(full_path):
            basename = os.path.basename(full_path).upper()
            # Only read files that look like license files, not METADATA etc.
            if any(basename.startswith(p.rstrip("*")) for p in LICENSE_FILE_PATTERNS):
                try:
                    with open(full_path, "r", errors="replace") as f:
                        return f.read().strip()
                except Exception:
                    continue
    return None

def find_go_license(name, version, gomodcache):
    """Find license text for a Go module from the module cache."""
    if not gomodcache:
        return None
    encoded = go_module_cache_encode(name)
    mod_dir = os.path.join(gomodcache, f"{encoded}@{version}")
    return find_license_file_in_dir(mod_dir)

def fetch_spdx_license_text(license_id):
    """Fetch canonical license text from spdx.org license-list-data (fallback)."""
    url = SPDX_LICENSE_TEXT_URL.format(license_id)
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            return resp.read().decode("utf-8").strip()
    except Exception:
        return None

def main():
    data = json.load(sys.stdin)
    artifacts = data.get("artifacts", [])
    scan_dir = os.environ.get("SCAN_DIR", ".")
    gomodcache = get_gomodcache()

    # Collect unique packages with their licenses and metadata
    packages = {}
    for a in artifacts:
        key = (a["name"], a.get("version", "unknown"))
        if key in packages:
            continue

        spdx_ids = set()
        lic_locations = []
        for lic in a.get("licenses", []):
            lic_locations.extend(lic.get("locations", []))
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
            packages[key] = {
                "spdx_ids": sorted(spdx_ids),
                "language": a.get("language", ""),
                "lic_locations": lic_locations,
            }

    if not packages:
        print("No packages with license information found.", file=sys.stderr)
        sys.exit(1)

    # Collect SPDX fallback texts (deduplicated)
    all_license_ids = set()
    for info in packages.values():
        all_license_ids.update(info["spdx_ids"])

    spdx_fallback_cache = {}
    cache_dir = os.path.join(os.environ.get("TMPDIR", "/tmp"), "spdx-license-cache")
    os.makedirs(cache_dir, exist_ok=True)

    # Generate NOTICE
    print("This software includes external packages and source code.")
    print("The applicable license information is listed below:")

    from_actual = 0
    from_spdx = 0
    no_text = 0

    for (name, version), info in sorted(packages.items()):
        spdx_ids = info["spdx_ids"]

        # Try to find actual license text from package sources
        actual_text = None

        # 1. Try syft-reported license file locations
        actual_text = find_license_from_syft_locations(info["lic_locations"], scan_dir)

        # 2. For Go packages, try the module cache
        if not actual_text and info["language"] == "go":
            actual_text = find_go_license(name, version, gomodcache)

        print()
        print("----")
        print()
        print(f"Package: {name}:{version}")
        for lid in spdx_ids:
            print(f"License: {lid}")

        if actual_text:
            print()
            print(actual_text)
            from_actual += 1
        else:
            # Fall back to canonical SPDX texts
            found_any = False
            for lid in spdx_ids:
                if lid not in spdx_fallback_cache:
                    cache_file = os.path.join(cache_dir, f"{lid}.txt")
                    if os.path.exists(cache_file):
                        with open(cache_file, "r") as f:
                            spdx_fallback_cache[lid] = f.read()
                    else:
                        text = fetch_spdx_license_text(lid)
                        if text:
                            spdx_fallback_cache[lid] = text
                            with open(cache_file, "w") as f:
                                f.write(text)
                        else:
                            spdx_fallback_cache[lid] = None

                text = spdx_fallback_cache.get(lid)
                if text:
                    print()
                    print(text)
                    found_any = True

            if found_any:
                from_spdx += 1
            else:
                no_text += 1

    print(
        f"License sources: {from_actual} from actual files, "
        f"{from_spdx} from SPDX fallback, "
        f"{no_text} without text",
        file=sys.stderr,
    )

main()
'
