#!/usr/bin/env bash
set -euo pipefail

python3 - <<'PY'
import importlib.util
spec = importlib.util.find_spec("huggingface_hub")
if spec is None:
    raise SystemExit("huggingface_hub is not installed")
import huggingface_hub
print(huggingface_hub.__version__)
PY

python3 -m huggingface_hub.commands.huggingface_cli --help >/dev/null
