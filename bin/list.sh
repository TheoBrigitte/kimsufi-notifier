#!/bin/bash
#
# Display available servers from OVH Eco (including Kimsufi) catalog
#
# Allowed country codes:
#   CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
#
# Usage:
#   COUNTRY=FR names.sh

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

source "${SCRIPT_DIR}/../config.env"

echo_stderr() {
    >&2 echo "$@"
}

OVH_URL="https://eu.api.ovh.com/v1/order/catalog/public/eco?ovhSubsidiary=${COUNTRY}"

# Fetch servers from OVH API
echo_stderr "> fetching servers in $COUNTRY"
DATA=$(curl -qSs "${OVH_URL}")

# Check for error: empty data, invalid json, or empty list
if test -z "$DATA" || ! echo "$DATA" | jq -e . &>/dev/null || echo "$DATA" | jq -e '.plans | length == 0' &>/dev/null; then
  echo "> failed to fetch data from $OVH_URL"
  exit 1
fi
echo_stderr "> fetched  servers"

# Get currency code
CURRENCY="$(echo "$DATA" | jq -r '.locale.currencyCode')"

# Print servers
echo "$DATA" | \
  jq -r '.plans[] | [ .planCode, .blobs.commercial.range, .invoiceName, (.pricings[] | select(.phase == 1) | select(.mode == "default") | .price/100000000) ] | @tsv' | \
  sort -k2,2 -k4n,4 -k3h,3 -b | \
  column -t -C "name=PlanCode" -C "name=Category" -C "name=Name" -C "name=Price ($CURRENCY)" -o '    '
