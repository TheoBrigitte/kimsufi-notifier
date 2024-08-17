#!/bin/bash
#
# Check OVH Eco (including Kimsufi) server availability
#
# Allowed datacenters:
#   bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw
#
# Usage:
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg check.sh
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg OPSGENIE_API_KEY=******** check.sh

set -eu

OPSGENIE_API_URL="https://api.opsgenie.com/v2/alerts"
OVH_URL="https://eu.api.ovh.com/v1/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}&datacenters=${DATACENTERS}"

echo_stderr() {
    >&2 echo "$@"
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
  exit 1
fi

# Print availability
AVAILABLE_DATACENTERS="$(echo "$DATA" | jq -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | join(",")')"
echo_stderr "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS"

# Exit here when OPSGENIE_API_KEY variable is not set
if [ -z ${OPSGENIE_API_KEY+x} ]; then
  exit 0
fi

# Send notification via OpsGenie
message="$PLAN_CODE is available\nhttps://eco.ovhcloud.com/fr/?display=list&range=kimsufi ."
echo_stderr "> sending notification"
curl -X POST "$OPSGENIE_API_URL" \
    -H "Content-Type: application/json" \
    -H "Authorization: GenieKey $OPSGENIE_API_KEY" \
    -d'{"message": "'"$message"'"}'
echo_stderr
echo_stderr "> notification sent"
