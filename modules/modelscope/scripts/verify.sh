#!/usr/bin/env bash
set -euo pipefail

python3 - <<'PY'
import importlib.util
spec = importlib.util.find_spec("modelscope")
if spec is None:
    raise SystemExit("modelscope is not installed")
import modelscope
print(modelscope.__version__)
PY

python3 - <<'PY'
import importlib.metadata as md

eps = md.entry_points()
group = eps.select(group="console_scripts") if hasattr(eps, "select") else eps.get("console_scripts", [])
names = {ep.name for ep in group}
if "modelscope" not in names:
    raise SystemExit(f"Expected modelscope console script, found: {sorted(names)}")
PY
