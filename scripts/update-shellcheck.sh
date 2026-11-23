#!/usr/bin/env bash
# Replace the bundled shellcheck binaries in internal/shellcheck/bin with the
# binaries of the version from internal/shellcheck/bin/VERSION.
#
# Usage:
#   scripts/update-shellcheck.sh [VERSION]
#
# If VERSION is omitted, the version is read from internal/shellcheck/bin/VERSION.
# Renovate bumps the VERSION file, so a renovate PR is applied by running this
# script, committing the changes and pushing to the renovate branch.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$REPO_ROOT/internal/shellcheck/bin"
VERSION_FILE="$BIN_DIR/VERSION"

if [ $# -ge 1 ]; then
    VERSION="$1"
else
    VERSION="$(tr -d '[:space:]' < "$VERSION_FILE")"
fi

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid version: $VERSION (expected e.g. 0.11.0)" >&2
    exit 1
fi

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

BASE_URL="https://github.com/koalaman/shellcheck/releases/download/v${VERSION}"

download_and_extract() {
    local asset="$1"

    curl --fail --location --show-error --silent -o "$WORK_DIR/$asset" "$BASE_URL/$asset"
    case "$asset" in
    *.tar.xz)
        mkdir -p "$WORK_DIR/$asset.out"
        tar -xf "$WORK_DIR/$asset" -C "$WORK_DIR/$asset.out"
        ;;
    *.zip)
        mkdir -p "$WORK_DIR/$asset.out"
        unzip -q -o "$WORK_DIR/$asset" -d "$WORK_DIR/$asset.out"
        ;;
    *)
        echo "Unsupported asset: $asset" >&2
        exit 1
        ;;
    esac
}

# Release asset suffix and repository file name per supported platform
platforms=(
    "darwin.aarch64.tar.xz darwin-aarch64"
    "darwin.x86_64.tar.xz darwin-x86_64"
    "linux.aarch64.tar.xz linux-aarch64"
    "linux.x86_64.tar.xz linux-x86_64"
    "zip windows.exe"
)

echo "Updating bundled shellcheck to $VERSION"

for platform in "${platforms[@]}"; do
    read -r suffix target <<<"$platform"

    asset="shellcheck-v${VERSION}.${suffix}"
    download_and_extract "$asset"

    if [ "$suffix" == "zip" ]; then
        binary="$WORK_DIR/$asset.out/shellcheck.exe"
    else
        binary="$WORK_DIR/$asset.out/shellcheck-v${VERSION}/shellcheck"
    fi

    cp "$binary" "$BIN_DIR/$target"
    chmod 755 "$BIN_DIR/$target"
    echo "  Updated $target"
done

# README.txt and LICENSE.txt ship with every release archive and document the
# build date plus the release commit.
cp "$WORK_DIR/shellcheck-v${VERSION}.zip.out/README.txt" "$BIN_DIR/README.txt"
cp "$WORK_DIR/shellcheck-v${VERSION}.zip.out/LICENSE.txt" "$BIN_DIR/LICENSE.txt"

echo "$VERSION" >"$VERSION_FILE"

echo "Done. Commit the changes in $BIN_DIR and push to the renovate branch."
