#!/usr/bin/env bash
# Bypass goenv's go.mod "go 1.26" check when 1.26 is not installed locally.
# Uses an installed goenv version + GOTOOLCHAIN=auto so the Go toolchain
# downloads 1.26 as needed.
set -euo pipefail

cd "$(dirname "$0")"

pick_goenv_version() {
  if [[ -n "${GOENV_VERSION:-}" ]]; then
    echo "$GOENV_VERSION"
    return
  fi
  if ! command -v goenv >/dev/null 2>&1; then
    return
  fi
  # Prefer newest installed 1.25.x, else any installed version (skip system).
  local versions
  versions="$(goenv versions --bare 2>/dev/null || true)"
  if [[ -z "$versions" ]]; then
    return
  fi
  local v
  v="$(echo "$versions" | grep -E '^1\.25\.' | sort -V | tail -1 || true)"
  if [[ -z "$v" ]]; then
    v="$(echo "$versions" | grep -v '^system$' | sort -V | tail -1 || true)"
  fi
  echo "$v"
}

GOENV_VER="$(pick_goenv_version || true)"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"

if [[ -n "${GOENV_VER}" ]]; then
  export GOENV_VERSION="$GOENV_VER"
  echo "run.sh: GOENV_VERSION=${GOENV_VERSION} GOTOOLCHAIN=${GOTOOLCHAIN}" >&2
fi

# With no args: run this package.
# Known go subcommands (build, mod, …) pass through: ./run.sh build -o /tmp/bwapi .
# Anything else is forwarded as program args: ./run.sh apikey → go run . apikey
if [[ $# -eq 0 ]]; then
  set -- run .
elif [[ "$1" != "run" && "$1" != "build" && "$1" != "mod" && "$1" != "test" &&
        "$1" != "env" && "$1" != "version" && "$1" != "get" && "$1" != "install" &&
        "$1" != "list" && "$1" != "fmt" && "$1" != "vet" && "$1" != "generate" &&
        "$1" != "clean" && "$1" != "work" && "$1" != "tool" && "$1" != "doc" &&
        "$1" != "fix" ]]; then
  set -- run . "$@"
fi

exec go "$@"
