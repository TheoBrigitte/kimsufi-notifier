#!/usr/bin/env bash
#
# Author: Théo Brigitte
# Date: 2025-10-26

#
# Usage: check.sh [options]
#
# Check OVH Eco (including Kimsufi) server availability
#
# Arguments
#   -p, --plan-code     Plan code to check (e.g. 24ska01)
#   -d, --datacenters   Comma-separated list of datacenters to check availability for (default all)
#                         Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)
#   -o, --option        Additional options to check for specific server options
#                         format key=value
#                         use --show-options to see available options
#       --show-options  Show available options for the plan code, requires --plan-code and --country
#       --country       Country code
#                         Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
#                         Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS
#                         Allowed values with -e ovh-us : US
#   -e, --endpoint      OVH API endpoint (default: ovh-eu)
#                         Allowed values: ovh-eu, ovh-ca, ovh-us
#       --verbose       Enable verbose mode to display detailed results, requires --country
#       --debug         Enable debug mode (default: false)
#   -h, --help          Display this help message
#
#   Arguments can also be set as environment variables see config.env.example
#   Command line arguments take precedence over environment variables
#
# Environment variables
#     DISCORD_WEBHOOK       Webhook URL to use for Discord notification service
#     GOTIFY_URL            URL to use for Gotify notification service
#     GOTIFY_TOKEN          token to use for Gotify notification service
#     GOTIFY_PRIORITY       prority for Gotify notification service
#     OPSGENIE_API_KEY      API key for OpsGenie to receive notifications
#     TELEGRAM_BOT_TOKEN    Bot token for Telegram to receive notifications
#     TELEGRAM_CHAT_ID      Chat ID for Telegram to receive notifications
#     HEALTHCHECKS_IO_UUID  UUID for healthchecks.io to ping after successful run
#
# Example:
#   check.sh --plan-code 24ska01
#   check.sh --plan-code 24ska01 --datacenters fr,gra,rbx,sbg

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  sed -Ene '/#\s?Usage:/,/^([^#]|$)/{p; /^([^#]|$)/q}' "$0" | sed -e '$d; s/#\s\?//'
}

notify_discord() {
  local message="$1"
  if [ -z ${DISCORD_WEBHOOK+x} ]; then
    return
  fi

  BODY="{\"content\": \"$message\"}"

  echo "> sending Discord notification"
  RESULT="$(curl -sSX POST -H "Content-Type: application/json" "$DISCORD_WEBHOOK" -d "$BODY")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .message 1>/dev/null; then
    echo_stderr "$RESULT"
    echo_stderr "> failed Discord notification"
  else
    echo "> sent Discord notification"
  fi
}

notify_gotify() {
  local message="$1"
  if [ -z ${GOTIFY_TOKEN+x} ]; then
    return
  fi

  if [ -z ${GOTIFY_URL+x} ]; then
    return
  fi

  if [ -z ${GOTIFY_PRIORITY+x} ]; then
    return
  fi

  echo "> sending Gotify notification"
  RESULT="$(curl -sSX POST "$GOTIFY_URL/message?token=$GOTIFY_TOKEN" \
      -F "title=OVH Availability" \
      -F "message=$message" \
      -F "priority=$GOTIFY_PRIORITY")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .error 1>/dev/null; then
    echo_stderr "$RESULT"
    echo_stderr "> failed Gotify notification"
  else
    echo "> sent Gotify notification"
  fi
}

notify_opsgenie() {
  local message="$1"
  if [ -z ${OPSGENIE_API_KEY+x} ]; then
    return
  fi

  echo "> sending OpsGenie notification"
  RESULT="$(curl -sSX POST "$OPSGENIE_API_URL" \
      -H "Content-Type: application/json" \
      -H "Authorization: GenieKey $OPSGENIE_API_KEY" \
      -d'{"message": "'"$message"'"}')"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e '.result | length > 0' 1>/dev/null; then
    echo "> sent    OpsGenie notification"
  else
    echo_stderr "$RESULT"
    echo_stderr "> failed  OpsGenie notification"
  fi
}

notify_telegram() {
  local message="$1"
  if [ -z ${TELEGRAM_BOT_TOKEN+x} ] || [ -z ${TELEGRAM_CHAT_ID+x} ]; then
    return
  fi

  TELEGRAM_WEBHOOK_URL="${TELEGRAM_API_URL}/bot${TELEGRAM_BOT_TOKEN}/sendMessage"

  echo "> sending Telegram notification"
  RESULT="$(curl -sSX POST \
    -d chat_id="${TELEGRAM_CHAT_ID}" \
    -d text="${message}" \
    -d parse_mode="HTML" \
    "${TELEGRAM_WEBHOOK_URL}")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .ok 1>/dev/null; then
    echo "> sent    Telegram notification"
  else
    echo_stderr "$RESULT"
    echo_stderr "> failed  Telegram notification"
  fi
}

join() {
  local d=${1-} f=${2-}
  if shift 2; then
    printf %s "$f" "${@/#/$d}"
  fi
}

get_catalog() {
  local country="$1"

  if [ -z "${country}" ]; then
    echo_stderr "Error: country is not set"
    return 1
  fi
  country="${country^^}"

  ovh_url="${OVH_API_ENDPOINTS["$ENDPOINT"]}/order/catalog/public/eco?ovhSubsidiary=${country}"
  data=$(curl -qSs "${ovh_url}")

  if test -z "$data" || ! echo "$data" | $JQ_BIN -e . 1>/dev/null || echo "$data" | $JQ_BIN -e '.plans | length == 0' 1>/dev/null; then
    echo_stderr "> failed to fetch data from $ovh_url"
    exit 2
  fi

  echo "$data"
}

# request makes an HTTP request to the OVH API
# Usage: request METHOD ENDPOINT [DATA] [OPTIONS]
request() {
  local method="$1"
  local endpoint="$2"
  local data="${3-}"
  if [ $# -lt 3 ]; then
    shift 2
  else
    shift 3
  fi

  if echo "$@" | grep -q -- '-v' || $DEBUG; then
    set -x
  fi
  result="$(curl -sX "${method}" "${OVH_URL}${endpoint}" \
    --header "Accept: application/json"\
    --header "Content-Type: application/json" \
    --data "${data}" \
    -w '%{stderr}%{http_code}' \
    "$@" 2>$HTTP_CODE_FILE)"
  set +x

  http_code=$(cat "$HTTP_CODE_FILE")
  if [ $http_code -lt 200 ] || [ $http_code -gt 299 ]; then
    echo_stderr "> error http_code=$http_code request=$method $OVH_URL$endpoint"
    echo_stderr "$result"
    return 1
  fi

  echo "$result"
  return 0
}

print_server_options() {
  local plan_code="$1"
  local catalog="$2"

  local plan_data="$(echo "$catalog" | $JQ_BIN -r '.plans[] | select(.planCode == "'"$plan_code"'")')"
  local products="$(echo "$catalog" | $JQ_BIN -r '.products')"

  output=""

  exec 6<<<$(echo "$plan_data" | $JQ_BIN -cr '.addonFamilies[]')
  while read <&6 addon; do
    mandatory="$(echo "$addon" | $JQ_BIN -r '.mandatory')"
    if [ "$mandatory" != "true" ]; then
      continue
    fi
    name="$(echo "$addon" | $JQ_BIN -r '.name')"
    if [ "$name" != "memory" ] && [ "$name" != "storage" ]; then
      continue
    fi

    default="$(echo "$addon" | $JQ_BIN -r '.default')"

    exec 7<<<$(echo "$addon" | $JQ_BIN -cr '.addons[]')
    while read <&7 value; do
      # cut last field
      real_value="$(echo "$value" | rev | cut -d'-' -f 2- | rev)"
      is_default=false
      if [[ "$value" == "$default" ]]; then
        is_default=true
      fi
      description="$(echo "$products" | $JQ_BIN -r '.[] | select(.name == "'"$real_value"'") | .description')"

      output+="$name=$real_value:$description:$is_default\n"
    done
  done
  exec 6<&-

  echo -e "$output" | column -t -s ':' -N "Option,Description,Default" -o '    '
}

print_verbose_availability() {
  local data="$1"

  output=""

  exec 6<<<$(echo "$data" | $JQ_BIN -cr '.[]')
  while read <&6 availability; do
    memory="$(echo "$availability" | $JQ_BIN -r '.memory')"
    storage="$(echo "$availability" | $JQ_BIN -r '.storage')"
    datacenters="$(echo "$availability" | $JQ_BIN -r '[.datacenters[] | select(.availability != "unavailable") | .datacenter] | unique | join(",")')"
    if [ -z "$datacenters" ]; then
      datacenters="unavailable"
    fi

    output+="  $memory:$storage:$datacenters\n"
  done
  exec 6<&-

  echo -e "$output" | column -t -s ':' -N "  Memory,Storage,Datacenters" -o '    '
}

main() {
  source "${SCRIPT_DIR}/common.sh"

  # Temporary file used to store HTTP reponse code
  HTTP_CODE_FILE="$(mktemp -t kimsufi-notifier.XXXXXX)"
  trap 'rm -f "$HTTP_CODE_FILE"' EXIT

  install_tools

  local options=()

  ARGS=$(getopt -o 'd:e:ho:p:v' --long 'country:,datacenters:,debug,endpoint:,help,option:,plan-code:,show-options,verbose' -- "$@")
  eval set -- "$ARGS"
  while true; do
    case "$1" in
      --country)
        COUNTRY="$2"
        shift 2
        continue
        ;;
      -d | --datacenters)
        DATACENTERS="$2"
        shift 2
        continue
        ;;
      --debug)
        DEBUG=true
        shift 1
        continue
        ;;
      -e | --endpoint)
        ENDPOINT="$2"
        shift 2
        continue
        ;;
      -h | --help)
        usage
        exit 0
        ;;
      -o | --option)
        echo "$2" | grep -q '=' || \
          exit_error "Error: invalid option '$2'"
        options+=("$2")
        shift 2
        continue
        ;;
      --show-options)
        SHOW_OPTIONS=true
        shift 1
        continue
        ;;
      -p | --plan-code)
        PLAN_CODE="$2"
        shift 2
        continue
        ;;
      -v | --verbose)
        VERBOSE=true
        shift 1
        continue
        ;;
      '--')
        shift
        break
        ;;
      *)
        echo_stderr 'Internal error!'
        exit 3
        ;;
    esac
  done

  if [ -z "${PLAN_CODE-}" ]; then
    echo_stderr "Error: PLAN_CODE is not set"
    echo_stderr
    usage
    exit 3
  fi

  if $SHOW_OPTIONS; then
    catalog=$(get_catalog "$COUNTRY")
    print_server_options "$PLAN_CODE" "$catalog"
    exit 0
  fi

  OVH_URL="${OVH_API_ENDPOINTS["$ENDPOINT"]}"
  endpoint="/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}"

  DATACENTERS_MESSAGE=""
  if [ -n "${DATACENTERS-}" ]; then
    endpoint="${endpoint}&datacenters=${DATACENTERS}"
    DATACENTERS_MESSAGE="$DATACENTERS datacenter(s)"
  else
    DATACENTERS_MESSAGE="all datacenters"
  fi

  if [ ${#options[@]} -gt 0 ]; then
    endpoint="${endpoint}&$(join '\&' "${options[@]}")"
  fi

  # Fetch availability from api
  echo "> checking $PLAN_CODE availability in $DATACENTERS_MESSAGE"
  if $DEBUG; then
    echo_stderr "> fetching data from $endpoint"
  fi

  DATA="$(request GET "${endpoint}")"

  if $DEBUG; then
    TMP_FILE="$(mktemp kimsufi-notifier.XXXXXX)"
    echo "$DATA" | tee "$TMP_FILE" 1>&2
    echo_stderr "> saved    data to   $TMP_FILE"
  fi

  # Check for error: empty data, invalid json, or empty list
  if test -z "$DATA" || ! echo "$DATA" | $JQ_BIN -e . 1>/dev/null || echo "$DATA" | $JQ_BIN -e '. | length == 0' 1>/dev/null; then
    echo_stderr "> failed to fetch data from $endpoint"
    exit 2
  fi

  # Ping healthchecks.io to ensure this script is running without errors
  if [ -n "${HEALTHCHECKS_IO_UUID-}" ]; then
    curl -sS -o /dev/null "${HEALTHCHECKS_IO_API_URL}/${HEALTHCHECKS_IO_UUID}"
  fi

  # Check for datacenters availability
  if ! echo "$DATA" | $JQ_BIN -e '.[].datacenters[] | select(.availability != "unavailable")' 1>/dev/null; then
    echo "> checked  $PLAN_CODE unavailable  in $DATACENTERS_MESSAGE"
    exit 4
  fi

  # Print availability
  AVAILABLE_DATACENTERS="$(echo "$DATA" | $JQ_BIN -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | unique | join(",")')"
  echo "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS datacenter(s)"
  if $VERBOSE; then
    print_verbose_availability "$DATA"
  fi

  # Send notifications
  message="$PLAN_CODE is available in $AVAILABLE_DATACENTERS datacenter(s), check https://eco.ovhcloud.com"
  notify_opsgenie "$message"
  notify_telegram "$message"
  notify_gotify "$message"
  notify_discord "$message"
}

main "$@"
