#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <service-file>"
  exit 1
fi

service_source="$1"

if [[ ! -f "${service_source}" ]]; then
  echo "Service file not found: ${service_source}"
  exit 1
fi

sudo_cmd=""
if [[ "$(id -u)" -ne 0 ]]; then
  sudo_cmd="sudo"
fi

${sudo_cmd} install -m 0644 "${service_source}" /etc/systemd/system/ollama.service
${sudo_cmd} systemctl daemon-reload
${sudo_cmd} systemctl enable --now ollama
