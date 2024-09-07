#!/bin/bash
#
# Check OVH Eco (including Kimsufi) server availability

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)
DEBUG=false

OVH_API_ENDPOINT="https://eu.api.ovh.com/v1"
OPSGENIE_API_URL="https://api.opsgenie.com/v2/alerts"
TELEGRAM_API_URL="https://api.telegram.org"
HEALTHCHECKS_IO_API_URL="https://hc-ping.com"

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  echo_stderr "Usage: PLAN_CODE=<plan code> DATACENTERS=<datacenters list> $0"
  echo_stderr "  Required:"
  echo_stderr "    PLAN_CODE             Plan code to check (e.g. 22sk010)"
  echo_stderr "  Optional:"
  echo_stderr "    DATACENTERS           Comma-separated list of datacenters"
  echo_stderr "                          Allowed values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw"
  echo_stderr "    OPSGENIE_API_KEY      API key for OpsGenie"
  echo_stderr "    TELEGRAM_BOT_TOKEN    Bot token for Telegram"
  echo_stderr "    TELEGRAM_CHAT_ID      Chat ID for Telegram"
  echo_stderr "    HEALTHCHECKS_IO_UUID  UUID for healthchecks.io"
  echo_stderr "    DEBUG                 Enable debug mode (default: false)"
  echo_stderr
  echo_stderr "Example:"
  echo_stderr "  PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg $0"
  echo_stderr "  PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg OPSGENIE_API_KEY=******* $0"
  echo_stderr "  PLAN_CODE=22sk010 DATACENTERS=fr,gra,rbx,sbg TELEGRAM_BOT_TOKEN=******** TELEGRAM_CHAT_ID=******** $0"
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

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e '.result | length > 0' &>/dev/null; then
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

  TELEGRAM_WEBHOOK_URL="${TELEGRAM_API_URL}/bot${TELEGRAM_BOT_TOKEN}/sendMessage"

  echo_stderr "> sending Telegram notification"
  RESULT="$(curl -sSX POST \
    -d chat_id="${TELEGRAM_CHAT_ID}" \
    -d text="${message}" \
    -d parse_mode="HTML" \
    "${TG_WEBHOOK_URL}")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .ok &>/dev/null; then
    echo_stderr "> sent    Telegram notification"
  else
    echo "$RESULT"
    echo_stderr "> failed  Telegram notification"
  fi
}

main() {
  source "${SCRIPT_DIR}/../config.env"
  source "${SCRIPT_DIR}/common.sh"

  ARGS=$(getopt -o 'de:hp:' --long 'datacenters:,debug,help,plan-code:' -- "$@")
  eval set -- "$ARGS"
  while true; do
    case "$1" in
      --datacenters)
        DATACENTERS="$2"
        shift 2
        continue
        ;;
      -d | --debug)
        DEBUG=true
        shift 1
        continue
        ;;
      -h | --help)
        usage
        exit 0
        ;;
      -p | --plan-code)
        PLAN_CODE="$2"
        shift 2
        continue
        ;;
      '--')
        shift
        break
        ;;
      *)
        echo_stderr 'Internal error!'
        exit 1
        ;;
    esac
  done

  if [ -z "${PLAN_CODE-}" ]; then
    echo_stderr "Error: PLAN_CODE is not set"
    echo_stderr
    usage
    exit 1
  fi

  OVH_URL="${OVH_API_ENDPOINT}/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}"

  DATACENTERS_MESSAGE=""
  if [ -n "${DATACENTERS-}" ]; then
    OVH_URL="${OVH_URL}&datacenters=${DATACENTERS}"
    DATACENTERS_MESSAGE="$DATACENTERS datacenter(s)"
  else
    DATACENTERS_MESSAGE="all datacenters"
  fi

  # Fetch availability from api
  echo_stderr "> checking $PLAN_CODE availability in $DATACENTERS_MESSAGE"
  if $DEBUG; then
    echo_stderr "> fetching data from $OVH_URL"
  fi

  DATA="$(curl -Ss "${OVH_URL}")"

  if $DEBUG; then
    TMP_FILE="$(mktemp kimsufi-notifier.XXXXXX)"
    echo "$DATA" | tee "$TMP_FILE"
    echo_stderr "> saved    data to   $TMP_FILE"
  fi

  # Check for error: empty data, invalid json, or empty list
  if test -z "$DATA" || ! echo "$DATA" | $JQ_BIN -e . &>/dev/null || echo "$DATA" | $JQ_BIN -e '. | length == 0' &>/dev/null; then
    echo "> failed to fetch data from $OVH_URL"
    exit 1
  fi

  # Ping healthchecks.io to ensure this script is running without errors
  if [ -n "${HEALTHCHECKS_IO_UUID-}" ]; then
    curl -sS -o /dev/null "${HEALTHCHECKS_IO_API_URL}/${HEALTHCHECKS_IO_UUID}"
  fi

  # Check for datacenters availability
  if ! echo "$DATA" | $JQ_BIN -e '.[].datacenters[] | select(.availability != "unavailable")' &>/dev/null; then
    echo_stderr "> checked  $PLAN_CODE unavailable  in $DATACENTERS_MESSAGE"
    exit 0
  fi

  # Print availability
  AVAILABLE_DATACENTERS="$(echo "$DATA" | $JQ_BIN -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | join(",")')"
  echo_stderr "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS"

  # Send notifications
  message="$PLAN_CODE is available https://eco.ovhcloud.com"
  notify_opsgenie "$message"
  notify_telegram "$message"
}

main "$@"
