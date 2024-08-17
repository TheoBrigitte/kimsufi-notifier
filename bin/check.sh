#!/bin/bash
#
# Check OVH Eco (including Kimsufi) server availability
#
# Allowed datacenters:
#   bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw
#
# Usage:
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg check.sh
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg OPSGENIE_API_KEY=******* check.sh
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg TELEGRAM_BOT_TOKEN=******** TELEGRAM_CHAT_ID=******** check.sh

set -eu

OPSGENIE_API_URL="https://api.opsgenie.com/v2/alerts"
OVH_URL="https://eu.api.ovh.com/v1/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}&datacenters=${DATACENTERS}"

echo_stderr() {
    >&2 echo "$@"
}

notify_opsgenie() {
  local message="$1"
  if [ -z ${OPSGENIE_API_KEY+x} ]; then
    return
  fi

  echo_stderr "> sending OpsGenie notification"
  RESULT="$(curl -sSX POST "$OPSGENIE_API_URL" \
      -H "Content-Type: application/json" \
      -H "Authorization: GenieKey $OPSGENIE_API_KEY" \
      -d'{"message": "'"$message"'"}')"

  if echo "$RESULT" | jq -e '.result | length > 0' &>/dev/null; then
    echo_stderr "> sent    OpsGenie notification"
  else
    echo "$RESULT"
    echo_stderr "> failed  OpsGenie notification"
  fi
}

notify_telegram() {
  local message="$1"
  if [ -z ${TELEGRAM_BOT_TOKEN+x} ] || [ -z ${TELEGRAM_CHAT_ID+x} ]; then
    return
  fi

  TG_WEBHOOK_URL="https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage"

  echo_stderr "> sending Telegram notification"
  RESULT="$(curl -sSX POST \
    -d chat_id="${TELEGRAM_CHAT_ID}" \
    -d text="${message}" \
    -d parse_mode="HTML" \
    "${TG_WEBHOOK_URL}")"

  if echo "$RESULT" | jq -e .ok &>/dev/null; then
    echo_stderr "> sent    Telegram notification"
  else
    echo "$RESULT"
    echo_stderr "> failed  Telegram notification"
  fi
}

# Fetch availability from api
echo_stderr "> checking $PLAN_CODE availability in $DATACENTERS"
DATA="$(curl -Ss "${OVH_URL}")"

# Check for error: empty data, invalid json, or empty list
if test -z "$DATA" || ! echo "$DATA" | jq -e . &>/dev/null || echo "$DATA" | jq -e '. | length == 0' &>/dev/null; then
  echo "> failed to fetch data from $OVH_URL"
  exit 1
fi

# Check for datacenters availability
if ! echo "$DATA" | jq -e '.[].datacenters[] | select(.availability != "unavailable")' &>/dev/null; then
  echo_stderr "> checked  $PLAN_CODE unavailable  in $DATACENTERS"
  exit 0
fi

# Print availability
AVAILABLE_DATACENTERS="$(echo "$DATA" | jq -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | join(",")')"
echo_stderr "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS"

# Send notifications
message="$PLAN_CODE is available https://eco.ovhcloud.com"
notify_opsgenie "$message"
notify_telegram "$message"
