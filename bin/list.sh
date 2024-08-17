#!/bin/bash
#
# Display available servers from OVH eco catalog
#
# Allowed country codes:
#   CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
#
# Usages:
#   COUNTRY=FR names.sh

set -eu

# Helper function - prints a message to stderr
echo_stderr() {
    >&2 echo "$@"
}

## required environement variables
_=$COUNTRY

OVH_URL="https://eu.api.ovh.com/v1/order/catalog/public/eco?ovhSubsidiary=${COUNTRY}"

echo_stderr "> fetching servers in $COUNTRY"
DATA=$(curl -qSs "${OVH_URL}")
if test -z "$DATA" || ! echo "$DATA" | jq -e . &>/dev/null || echo "$DATA" | jq -e '.plans | length == 0' &>/dev/null; then
  echo "> failed to fetch data from $OVH_URL"
  exit 1
fi
echo_stderr "> fetched  servers"

CURRENCY="$(echo "$DATA" | jq -r '.locale.currencyCode')"
echo "$DATA" | \
  jq -r '.plans[] | [.planCode,.blobs.commercial.range,.invoiceName,(.pricings[]|select(.phase == 1)|select(.mode == "default")|.price/100000000)] | @tsv' | \
  sort -k2,2 -k4n,4 -k3h,3 -b | \
  column -t -C "name=PlanCode" -C "name=Category" -C "name=Name" -C "name=Price ($CURRENCY)" -o '    '
