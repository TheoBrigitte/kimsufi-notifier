#!/usr/bin/env bash
#
# Author: Théo Brigitte
# Date: 2025-10-26

# Usage: list.sh [options]
#
# List servers from OVH Eco (including Kimsufi) catalog
#
# Arguments
#       --category   Server category (default all)
#                      Allowed values: kimsufi, soyoustart, rise, uncategorized
#   -c, --country    Country code (required)
#                      Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
#                      Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS
#                      Allowed values with -e ovh-us : US
#   -e, --endpoint   OVH API endpoint (default: ovh-eu)
#                      Allowed values: ovh-eu, ovh-ca, ovh-us
#       --debug      Enable debug mode (default: false)
#   -h, --help       Display this help message
#
#   Arguments can also be set as environment variables see config.env.example
#   Command line arguments take precedence over environment variables
#
# Example:
#     list.sh --country FR
#     list.sh --country FR --category kimsufi

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

DEBUG=false
ENDPOINT="ovh-eu"

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  sed -Ene '/#\s?Usage:/,/^([^#]|$)/{p; /^([^#]|$)/q}' "$0" | sed -e '$d; s/#\s\?//'
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
        ENDPOINT="$2"
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
        exit 3
        ;;
    esac
  done

  if [ -z "${COUNTRY-}" ]; then
    echo_stderr "Error: COUNTRY is not set"
    echo_stderr
    usage
    exit 3
  fi
  COUNTRY="${COUNTRY^^}"

  OVH_URL="${OVH_API_ENDPOINTS["$ENDPOINT"]}/order/catalog/public/eco?ovhSubsidiary=${COUNTRY}"

  # Fetch servers from OVH API
  echo "> fetching servers in $COUNTRY"
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
  if test -z "$DATA" || ! echo "$DATA" | $JQ_BIN -e . 1>/dev/null || echo "$DATA" | $JQ_BIN -e '.plans | length == 0' 1>/dev/null; then
    echo_stderr "> failed to fetch data from $OVH_URL"
    exit 2
  fi
  echo "> fetched  servers"

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
    column -s $'\t' -t -N "PlanCode,Category,Name,Price ($CURRENCY)" -o '    '
}

main "$@"
