#!/usr/bin/env bash
# Generate plugins from specs directory using assistantkit
#
# This script regenerates all plugin outputs from the canonical specs in specs/.
# It replaces the old plugins/generate/main.go approach.
#
# Usage:
#   ./scripts/generate-plugins.sh
#
# Requirements:
#   - assistantkit CLI (or go run from assistantkit source)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ASSISTANTKIT_SRC="${PROJECT_ROOT}/../assistantkit"

if command -v assistantkit &> /dev/null; then
    assistantkit generate all \
        --specs="${PROJECT_ROOT}/specs" \
        --target=local \
        --output="${PROJECT_ROOT}"
elif [[ -d "$ASSISTANTKIT_SRC" ]]; then
    (cd "$ASSISTANTKIT_SRC" && go run ./cmd/assistantkit generate all \
        --specs="${PROJECT_ROOT}/specs" \
        --target=local \
        --output="${PROJECT_ROOT}")
else
    echo "Error: assistantkit not found. Install it or clone to ../assistantkit"
    exit 1
fi
