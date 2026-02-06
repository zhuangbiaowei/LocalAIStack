#!/usr/bin/env bash
set -euo pipefail

VENV_DIR="${VLLM_VENV_DIR:-$HOME/.localaistack/venv/vllm}"
SOURCE_VENV="${VLLM_SOURCE_DIR:-$HOME/vllm}/.venv"

if [[ ! -x "$VENV_DIR/bin/python" && -x "$SOURCE_VENV/bin/python" ]]; then
  VENV_DIR="$SOURCE_VENV"
fi

VENV_PY="$VENV_DIR/bin/python"

if [[ ! -x "$VENV_PY" ]]; then
  echo "vLLM venv not found at $VENV_DIR" >&2
  exit 1
fi

"$VENV_PY" - <<'PY'
import vllm
print(vllm.__version__)
PY
