#!/bin/bash
#
# Check OVH Eco (including Kimsufi) server availability

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)
DEBUG=false
SHOW_OPTIONS=false
VERBOSE=false

ENDPOINT="ovh-eu"
OPSGENIE_API_URL="https://api.opsgenie.com/v2/alerts"
TELEGRAM_API_URL="https://api.telegram.org"
HEALTHCHECKS_IO_API_URL="https://hc-ping.com"

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  bin_name=$(basename "$0")
  echo_stderr "Usage: $bin_name"
  echo_stderr
  echo_stderr "Check OVH Eco (including Kimsufi) server availability"
  echo_stderr
  echo_stderr "Arguments"
  echo_stderr "  -p, --plan-code     Plan code to check (e.g. 24ska01)"
  echo_stderr "  -d, --datacenters   Comma-separated list of datacenters to check availability for (default all)"
  echo_stderr "                        Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)"
  echo_stderr "  -o, --option        Additional options to check for specific server options"
  echo_stderr "                        format key=value"
  echo_stderr "                        use --show-options to see available options"
  echo_stderr "      --show-options  Show available options for the plan code, requires --plan-code and --country"
  echo_stderr "      --country       Country code"
  echo_stderr "                        Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN"
  echo_stderr "                        Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS"
  echo_stderr "                        Allowed values with -e ovh-us : US"
  echo_stderr "  -e, --endpoint      OVH API endpoint (default: ovh-eu)"
  echo_stderr "                        Allowed values: ovh-eu, ovh-ca, ovh-us"
  echo_stderr "      --verbose       Enable verbose mode to display detailed results, requires --country"
  echo_stderr "      --debug         Enable debug mode (default: false)"
  echo_stderr "  -h, --help          Display this help message"
  echo_stderr
  echo_stderr "  Arguments can also be set as environment variables see config.env.example"
  echo_stderr "  Command line arguments take precedence over environment variables"
  echo_stderr
  echo_stderr "Environment variables"
  echo_stderr "    DISCORD_WEBHOOK       Webhook URL to use for Discord notification service"
  echo_stderr "    GOTIFY_URL            URL to use for Gotify notification service"
  echo_stderr "    GOTIFY_TOKEN          token to use for Gotify notification service"
  echo_stderr "    GOTIFY_PRIORITY       prority for Gotify notification service"
  echo_stderr "    OPSGENIE_API_KEY      API key for OpsGenie to receive notifications"
  echo_stderr "    TELEGRAM_BOT_TOKEN    Bot token for Telegram to receive notifications"
  echo_stderr "    TELEGRAM_CHAT_ID      Chat ID for Telegram to receive notifications"
  echo_stderr "    HEALTHCHECKS_IO_UUID  UUID for healthchecks.io to ping after successful run"
  echo_stderr
  echo_stderr "Example:"
  echo_stderr "  $bin_name --plan-code 24ska01"
  echo_stderr "  $bin_name --plan-code 24ska01 --datacenters fr,gra,rbx,sbg"
}

notify_discord() {
  local message="$1"
  if [ -z ${DISCORD_WEBHOOK+x} ]; then
    return
  fi

  BODY="{\"content\": \"$message\"}"

  echo_stderr "> sending Discord notification"
  RESULT="$(curl -sSX POST -H "Content-Type: application/json" "$DISCORD_WEBHOOK" -d "$BODY")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .message &>/dev/null; then
    echo "$RESULT"
    echo_stderr "> failed Discord notification"
  else
    echo_stderr "> sent Discord notification"
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

  echo_stderr "> sending Gotify notification"
  RESULT="$(curl -sSX POST "$GOTIFY_URL/message?token=$GOTIFY_TOKEN" \
      -F "title=OVH Availability" \
      -F "message=$message" \
      -F "priority=$GOTIFY_PRIORITY")"

  if $DEBUG; then
    echo_stderr "$RESULT"
  fi

  if echo "$RESULT" | $JQ_BIN -e .error &>/dev/null; then
    echo "$RESULT"
    echo_stderr "> failed Gotify notification"
  else
    echo_stderr "> sent Gotify notification"
  fi
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
    "${TELEGRAM_WEBHOOK_URL}")"

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

  if test -z "$data" || ! echo "$data" | $JQ_BIN -e . &>/dev/null || echo "$data" | $JQ_BIN -e '.plans | length == 0' &>/dev/null; then
    echo_stderr "> failed to fetch data from $ovh_url"
    exit 1
  fi

  echo "$data"
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
  source "${SCRIPT_DIR}/../config.env"
  source "${SCRIPT_DIR}/common.sh"

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

  if $SHOW_OPTIONS; then
    catalog=$(get_catalog "$COUNTRY")
    print_server_options "$PLAN_CODE" "$catalog"
    exit 0
  fi

  OVH_URL="${OVH_API_ENDPOINTS["$ENDPOINT"]}/dedicated/server/datacenter/availabilities?planCode=${PLAN_CODE}"

  DATACENTERS_MESSAGE=""
  if [ -n "${DATACENTERS-}" ]; then
    OVH_URL="${OVH_URL}&datacenters=${DATACENTERS}"
    DATACENTERS_MESSAGE="$DATACENTERS datacenter(s)"
  else
    DATACENTERS_MESSAGE="all datacenters"
  fi

  if [ ${#options[@]} -gt 0 ]; then
    OVH_URL="${OVH_URL}&$(join '\&' "${options[@]}")"
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
    exit 1
  fi

  # Print availability
  AVAILABLE_DATACENTERS="$(echo "$DATA" | $JQ_BIN -r '[.[].datacenters[] | select(.availability != "unavailable") | .datacenter] | unique | join(",")')"
  echo_stderr "> checked  $PLAN_CODE available    in $AVAILABLE_DATACENTERS datacenter(s)"
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
