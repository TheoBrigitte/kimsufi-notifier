#!/bin/bash
#
# Display available servers from OVH Eco (including Kimsufi) catalog

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

source "${SCRIPT_DIR}/../config.env"
source "${SCRIPT_DIR}/common.sh"

echo_stderr() {
    >&2 echo "$@"
}

usage() {
  echo_stderr "Usage: COUNTRY=XX $0"
  echo_stderr "  Required:"
  echo_stderr "    COUNTRY             Country code"
  echo_stderr "                        Allowed values: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN"
  echo_stderr
  echo_stderr "Optional:"
  echo_stderr "    DEBUG               Enable debug mode (default: false)"
  echo_stderr
  echo_stderr "Example:"
  echo_stderr "  COUNTRY=FR $0"
}

if [ -z "${COUNTRY-}" ]; then
  echo_stderr "Error: COUNTRY is not set"
  echo_stderr
  usage
  exit 1
fi

OVH_URL="https://eu.api.ovh.com/v1/order/catalog/public/eco?ovhSubsidiary=${COUNTRY}"

DEBUG=${DEBUG:-false}

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

# Print servers
echo "$DATA" | \
  $JQ_BIN -r '.plans[] | [ .planCode, .blobs.commercial.range, .invoiceName, (.pricings[] | select(.phase == 1) | select(.mode == "default") | .price/100000000) ] | @tsv' | \
  sort -k2,2 -k4n,4 -b -t $'\t' | \
  column -s $'\t' -t -C "name=PlanCode" -C "name=Category" -C "name=Name" -C "name=Price ($CURRENCY)" -o '    '
