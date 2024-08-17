#!/bin/bash
#
# Check Kimsufi server availability
#
# Usages:
# 	PLAN_CODE=22sk010 DATACENTERS=fr OPSGENIE_API_KEY=******** check.sh
# 	PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg NO_NOTIFICATION=true check.sh

set -eu

# Helper function - prints a message to stderr
echo_stderr() {
    >&2 echo "$@"
}

## required environement variables
_=$PLAN_CODE
_=$DATACENTERS

OPSGENIE_API_URL="https://api.opsgenie.com/v2/alerts"
OVH_URL="https://eu.api.ovh.com/v1/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}&datacenters=${DATACENTERS}"

# check availability from api
echo_stderr "> checking $PLAN_CODE availability in $DATACENTERS"
DATA="$(curl -Ss "${OVH_URL}")"
#DATA="$(echo bob)"
if test -z "$DATA" || ! echo "$DATA" | jq -e . &>/dev/null || echo "$DATA" | jq -e '. | length == 0' &>/dev/null; then
  echo "> failed to fetch data from $OVH_URL"
  exit 1
fi

if ! echo "$DATA" | jq -e '.[].datacenters[] | select(.availability != "unavailable")' &>/dev/null; then
  echo_stderr "> checked  $PLAN_CODE unavailable  in $DATACENTERS"
  exit 1
fi

AVAILABLE_DATACENTERS="$(echo "$DATA" | jq -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | join(",")')"
echo_stderr "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS"

# stop when NO_NOTIFICATION variable is set
if $NO_NOTIFICATION; then
  exit 0
fi

_=$OPSGENIE_API_KEY

# send notification
message="$PLAN_CODE is available\nhttps://eco.ovhcloud.com/fr/?display=list&range=kimsufi ."
echo_stderr "> sending notification"
curl -X POST "$OPSGENIE_API_URL" \
    -H "Content-Type: application/json" \
    -H "Authorization: GenieKey $OPSGENIE_API_KEY" \
    -d'{"message": "'"$message"'"}'
echo_stderr
echo_stderr "> notification sent"
