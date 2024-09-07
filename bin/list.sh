#!/bin/bash
#
# Display available servers from OVH Eco (including Kimsufi) catalog

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

DEBUG=false

OVH_API_ENDPOINT="ovh-eu"

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  bin_name=$(basename "$0")
  echo_stderr "Usage: $bin_name"
  echo_stderr
  echo_stderr "List servers from OVH Eco (including Kimsufi) catalog"
  echo_stderr
  echo_stderr "Arguments"
  echo_stderr "  -c, --country    Country code (required)"
  echo_stderr "                     Allowed values for ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN"
  echo_stderr "                     Allowed values for ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS"
  echo_stderr "                     Allowed values for ovh-us: US"
  echo_stderr "  --category       Server category (default all)"
  echo_stderr "                     Allowed values: kimsufi, soyoustart, rise, uncategorized"
  echo_stderr "  -e, --endpoint   OVH API endpoint (default: ovh-eu)"
  echo_stderr "                     Allowed values: ovh-eu, ovh-ca, ovh-us"
  echo_stderr "  -d, --debug      Enable debug mode (default: false)"
  echo_stderr "  -h, --help       Display this help message"
  echo_stderr
  echo_stderr "  Arguments can also be set as environment variables see config.env.example"
  echo_stderr "  Command line arguments take precedence over environment variables"
  echo_stderr
  echo_stderr "Example:"
  echo_stderr "    $bin_name --country FR"
}

main() {
  source "${SCRIPT_DIR}/../config.env"
  source "${SCRIPT_DIR}/common.sh"

  install_tools

  ARGS=$(getopt -o 'c:e:h' --long 'category:,country:,debug,endpoint:,help' -- "$@")
  eval set -- "$ARGS"
  while true; do
    case "$1" in
      --category)
        CATEGORY="$2"
        shift 2
        continue
        ;;
      -c | --country)
        COUNTRY="$2"
        shift 2
        continue
        ;;
      --debug)
        DEBUG=true
        shift 1
        continue
        ;;
      -e | --endpoint)
        OVH_API_ENDPOINT="$2"
        shift 2
        continue
        ;;
      -h | --help)
        usage
        exit 0
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

  if [ -z "${COUNTRY-}" ]; then
    echo_stderr "Error: COUNTRY is not set"
    echo_stderr
    usage
    exit 1
  fi
  COUNTRY="${COUNTRY^^}"

  OVH_URL="${OVH_API_ENDPOINTS["$OVH_API_ENDPOINT"]}/order/catalog/public/eco?ovhSubsidiary=${COUNTRY}"

  # Fetch servers from OVH API
  echo_stderr "> fetching servers in $COUNTRY"
  if $DEBUG; then
    echo_stderr "> fetching data from $OVH_URL"
  fi

  DATA=$(curl -qSs "${OVH_URL}")

  if $DEBUG; then
    TMP_FILE="$(mktemp kimsufi-notifier.XXXXXX)"
    echo "$DATA" > "$TMP_FILE"
    echo_stderr "> saved    data to   $TMP_FILE"
  fi

  # Check for error: empty data, invalid json, or empty list
  if test -z "$DATA" || ! echo "$DATA" | $JQ_BIN -e . &>/dev/null || echo "$DATA" | $JQ_BIN -e '.plans | length == 0' &>/dev/null; then
    echo "> failed to fetch data from $OVH_URL"
    exit 1
  fi
  echo_stderr "> fetched  servers"

  # Get currency code
  CURRENCY="$(echo "$DATA" | $JQ_BIN -r '.locale.currencyCode')"

  # Filter by category
  if [ -n "${CATEGORY-}" ]; then
    if [ "$CATEGORY" == "uncategorized" ]; then
      CATEGORY=null
    else
      CATEGORY="\"$CATEGORY\""
    fi
    category_filter='| select(.blobs.commercial.range == '"$CATEGORY"')'
  else
    category_filter=''
  fi

  # Print servers
  echo "$DATA" | \
    $JQ_BIN -r '.plans[] '"$category_filter"' | [ .planCode, .blobs.commercial.range, .invoiceName, (.pricings[] | select(.phase == 1) | select(.mode == "default") | .price/100000000) ] | @tsv' | \
    sort -k2,2 -k4n,4 -b -t $'\t' | \
    column -s $'\t' -t -C "name=PlanCode" -C "name=Category" -C "name=Name" -C "name=Price ($CURRENCY)" -o '    '
}

main "$@"
