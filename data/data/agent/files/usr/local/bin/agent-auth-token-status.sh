#!/bin/bash

set -e

# shellcheck disable=SC1091
source "issue_status.sh"

status_issue="65_token"

export TZ=UTC
check_token_expiry() {
  expiry_epoch=$(date -d "${AUTH_TOKEN_EXPIRY}" +%s)
  current_epoch=$(date +%s)

  if [ "$current_epoch" -gt "$expiry_epoch" ]; then
    printf '\\e{lightred}The authentication token has expired. Please generate a new ISO using the "oc adm node-image create" command, then reboot the node.\\e{reset}'| set_issue "${status_issue}"
    exit 1
  fi
}

while true; do
  check_token_expiry
  sleep 5
done
