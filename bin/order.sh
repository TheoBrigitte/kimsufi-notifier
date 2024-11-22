#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

# Default values
DEBUG=false
DRY_RUN=false
ENDPOINT="ovh-eu"
PRICE_DURATION="P1M"
PRICE_MODE="default"
QUANTITY=1

echo_stderr() {
    >&2 echo "$@"
}

# Helper function - prints an error message and exits
exit_error() {
    echo_stderr "Error: $1"
    exit 1
}

usage() {
  bin_name=$(basename "$0")
  echo_stderr "Usage: $bin_name"
  echo_stderr
  echo_stderr "Place an order for a servers from OVH Eco (including Kimsufi) catalog"
  echo_stderr
  echo_stderr "Arguments"
  echo_stderr "  -c, --country     Country code (required)"
  echo_stderr "                      Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN"
  echo_stderr "                      Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS"
  echo_stderr "                      Allowed values with -e ovh-us : US"
  echo_stderr "  --datacenter      Datacenter code (default: from config when only one is set)"
  echo_stderr "                      Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)"
  echo_stderr "  -e, --endpoint    OVH API endpoint (default: $ENDPOINT)"
  echo_stderr "                      Allowed values: ovh-eu, ovh-ca, ovh-us"
  echo_stderr "  -i, --item-configuration"
  echo_stderr "                      Item configuration in the form 'label=value'"
  echo_stderr "  -d, --debug       Enable debug mode (default: $DEBUG)"
  echo_stderr "      --dry-run     Do not create the order, only configure the cart (default: $DRY_RUN)"
  echo_stderr "  -h, --help        Display this help message"
  echo_stderr "  -q, --quantity    Quantity of items to order (default: $QUANTITY)"
  echo_stderr "  --price-mode      Billing price type (default: $PRICE_MODE)"
  echo_stderr "  --price-duration  Billing duration (default: $PRICE_DURATION)"
  echo_stderr
  echo_stderr "  Arguments can also be set as environment variables see config.env.example"
  echo_stderr "  Command line arguments take precedence over environment variables"
  echo_stderr
  echo_stderr "Example:"
  echo_stderr "    $bin_name"
  echo_stderr "    $bin_name --item-configuration region=europe"
  echo_stderr "    $bin_name --item-configuration region=europe --datacenter fra"
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
  curl -sX "${method}" "${OVH_URL}${endpoint}" \
    --header "Accept: application/json"\
    --header "Content-Type: application/json" \
    --data "${data}" \
    -w '\n%output{'$HTTP_CODE_FILE'}%{http_code}' \
    "$@"
  set +x

  http_code=$(cat "$HTTP_CODE_FILE")
  if [ $http_code -lt 200 ] || [ $http_code -gt 299 ]; then
    echo_stderr "> error http_code=$http_code request=$method $OVH_URL$endpoint"
    return 1
  fi

  return 0
}

# request_auth makes an authenticated HTTP request to the OVH API
# Usage: request_auth METHOD ENDPOINT [DATA]
request_auth() {
  local method="$1"
  local endpoint="$2"
  local data="${3-}"
  if [ $# -lt 3 ]; then
    shift 2
  else
    shift 3
  fi

  local timestamp="$(date +%s)"
  local sig_key="${APPLICATION_SECRET}+${CONSUMER_KEY}+${method}+${OVH_URL}${endpoint}+${data}+${timestamp}"
  local signature=$(echo "\$1\$$(echo -n "${sig_key}" | sha1sum - | awk '{print $1}')")

  request "${method}" "${endpoint}" "${data}" \
    --header "X-Ovh-Application: ${APPLICATION_KEY}" \
    --header "X-Ovh-Consumer: ${CONSUMER_KEY}" \
    --header "X-Ovh-Timestamp: ${timestamp}" \
    --header "X-Ovh-Signature: ${signature}" \
    "$@"
}

# item_auto_configuration automatically configures an item with required configuration having only one allowed value
item_auto_configuration() {
  local cart_id="$1"
  local item_id="$2"

  exec 6<<<$(request GET "/order/cart/${cart_id}/item/${item_id}/requiredConfiguration" | $JQ_BIN -cr '.[]|select((.required==true) or (.label=="dedicated_datacenter"))|select(.allowedValues|length == 1)')

  local labels=()
  while read <&6 configuration; do
    label="$(echo "$configuration" | $JQ_BIN -r .label)"
    value="$(echo "$configuration" | $JQ_BIN -r .allowedValues[0])"
    echo_stderr "> item auto-configuration $label=$value"
    request POST "/order/cart/${cart_id}/item/${item_id}/configuration" '{"label":"'"$label"'","value":"'"$value"'"}'
    labels+=("$label")
  done

  exec 6<&-

  echo "${labels[@]}"
}

# item_user_configuration configures an item with configuration passed as arguments
item_user_configuration() {
  local cart_id="$1"
  local item_id="$2"
  shift 2
  local configurations=("$@")

  local labels=()
  for configuration in ${configurations[@]}; do
    label="$(echo "$configuration" | cut -d= -f1)"
    value="$(echo "$configuration" | cut -d= -f2)"
    echo_stderr "> item user-configuration $label=$value"
    request POST "/order/cart/${cart_id}/item/${item_id}/configuration" '{"label":"'"$label"'","value":"'"$value"'"}'
    labels+=("$label")
  done

  echo "${labels[@]}"
}

# item_manual_configuration ask user to manually configures item with remaining required configuration
item_manual_configuration() {
  local cart_id="$1"
  local item_id="$2"
  shift 2
  local labels_configured=("$@")

  exec 6<<<$(request GET "/order/cart/${cart_id}/item/${item_id}/requiredConfiguration" | $JQ_BIN -cr '.[]|select((.required==true) or (.label=="dedicated_datacenter"))')

  while read <&6 configuration; do
    label="$(echo "$configuration" | $JQ_BIN -r .label)"
    if [[ ${labels_configured[@]} =~ $label ]]; then
      continue
    fi
    echo_stderr "> item configuration, select a value for $label"

    i=0
    for value in $(echo "$configuration" | $JQ_BIN -r '.allowedValues[]'); do
      echo_stderr "> $i. $value"
      i=$((i+1))
    done
    read -p "> Choice: " index
    value="$(echo "$configuration" | $JQ_BIN -r .allowedValues[$index])"
    echo_stderr "> item manual-configuration $label=$value"
    request POST "/order/cart/${cart_id}/item/${item_id}/configuration" '{"label":"'"$label"'","value":"'"$value"'"}'
  done

  exec 6<&-
}

# item_option_configuration configures item with mandatory options, choosing the cheapest available
item_option_configuration() {
  local cart_id="$1"
  local item_id=$2
  local plan_code="$3"
  local price_mode="$4"
  local price_duration="$5"

  exec 6<<<$(request GET "/order/cart/${cart_id}/eco/options?planCode=${plan_code}" | $JQ_BIN -cr '.[]')

  declare -A familyPlanCode
  declare -A familyPrices
  while read <&6 option; do
    mandatory="$(echo "$option" | $JQ_BIN -r .mandatory)"
    if [ "$mandatory" != "true" ]; then
      continue
    fi
    family="$(echo "$option" | $JQ_BIN -r .family)"
    code="$(echo "$option" | $JQ_BIN -r .planCode)"
    price="$(echo "$option" | $JQ_BIN -r '.prices[]|select((.pricingMode == "'"$price_mode"'") and (.duration == "'"$price_duration"'"))|.priceInUcents')"
    if ! [ "${familyPrices[$family]+x}" ]; then
      familyPrices[$family]=$price
      familyPlanCode[$family]="$code"
    elif [ "$price" -lt "${familyPrices[$family]}" ]; then
      familyPrices[$family]=$price
      familyPlanCode[$family]="$code"
    fi
  done

  exec 6<&-

  for option in "${familyPlanCode[@]}"; do
    echo_stderr "> item option $option"
    result="$(request POST "/order/cart/${cart_id}/eco/options" '{"quantity": 1, "duration": "'"$price_duration"'", "pricingMode":"'"$price_mode"'", "planCode":"'"$option"'", "itemId": '$item_id'}')"
    $DEBUG && echo "$result" $JQ_BIN -cr .
  done
}

main() {
  # Load configuration and common tools
  source "${SCRIPT_DIR}/../config.env"
  source "${SCRIPT_DIR}/common.sh"

  # Temporary file used to store HTTP reponse code
  HTTP_CODE_FILE="$(mktemp -t kimsufi-notifier.XXXXXX)"
  trap 'rm -f "$HTTP_CODE_FILE"' EXIT

  # Use configured dataceter if only one is set
  DATACENTER=""
  if [ -n "${DATACENTERS-}" ] && echo "$DATACENTERS"|grep -vq ,; then
    DATACENTER="$DATACENTERS"
  fi

  install_tools

  local item_configurations=()

  ARGS=$(getopt -o 'c:d:e:hi:p:q:' --long 'country:,datacenter:,item-configuration:,debug,dry-run,endpoint:,help,quantity:,plan-code:,price-duration:,price-mode:' -- "$@")
  eval set -- "$ARGS"
  while true; do
    case "$1" in
      -c | --country)
        COUNTRY="$2"
        shift 2
        continue
        ;;
      -d | --datacenter)
        DATACENTER="$2"
        shift 2
        continue
        ;;
      --debug)
        DEBUG=true
        shift 1
        continue
        ;;
      --dry-run)
        DRY_RUN=true
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
      -i | --item-configuration)
        echo "$2" | grep -q '=' || \
          exit_error "Error: invalid item configuration '$2'"
        item_configurations+=("$2")
        shift 2
        continue
        ;;
      -q | --quantity)
        QUANTITY="$2"
        shift 2
        continue
        ;;
      -p | --plan-code)
        PLAN_CODE="$2"
        shift 2
        continue
        ;;
      --price-mode)
        PRICE_MODE="$2"
        shift 2
        continue
        ;;
      --price-duration)
        PRICE_DURATION="$2"
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

  if [ -z "${COUNTRY-}" ]; then
    echo_stderr "Error: COUNTRY is not set"
    echo_stderr
    usage
    exit 1
  fi
  COUNTRY="${COUNTRY^^}"

  if [ -n "${DATACENTER-}" ]; then
    item_configurations+=("dedicated_datacenter=$DATACENTER")
  fi

  # OVH API endpoint
  OVH_URL="${OVH_API_ENDPOINTS["$ENDPOINT"]}"

  # Create cart
  expire="$(date --iso-8601=seconds --date tomorrow)"
  cart="$(request POST "/order/cart" '{"description":"kimsufi-notifier","expire":"'"$expire"'","ovhSubsidiary":"'"$COUNTRY"'"}')"
  $DEBUG && echo "$cart" | $JQ_BIN -cr .

  cart_id="$(echo "$cart" | $JQ_BIN -r .cartId)"
  if [ -z "$cart_id" ]; then
    echo "cart_id is empty"
    exit 1
  fi
  echo "> cart created id=$cart_id"

  # Add item to cart
  cart_updated="$(request POST "/order/cart/${cart_id}/eco" '{"planCode":"'"${PLAN_CODE}"'","quantity": '${QUANTITY}', "pricingMode":"'"${PRICE_MODE}"'","duration":"'${PRICE_DURATION}'"}')"
  $DEBUG && echo "$cart_updated" | $JQ_BIN -cr .

  item_id="$(echo "$cart_updated" | $JQ_BIN -r .itemId)"
  if [ -z "$item_id" ]; then
    echo "item_id is empty"
    exit 1
  fi
  echo "> cart updated with item id=$item_id"

  # Configure item
  labels_auto_configured="$(item_auto_configuration "$cart_id" "$item_id")"
  labels_user_configured="$(item_user_configuration "$cart_id" "$item_id" "${item_configurations[@]}")"
  labels_configured=( "${labels_auto_configured[@]}" "${labels_user_configured[@]}" )
  item_manual_configuration "$cart_id" "$item_id" "${labels_configured[@]}"

  # Configure eco options
  item_option_configuration "$cart_id" $item_id "$PLAN_CODE" "$PRICE_MODE" "$PRICE_DURATION"

  if $DRY_RUN; then
    echo_stderr "> dry-run enabled, skipping order completion"
    return
  fi

  # Assign cart to account
  request_auth POST "/order/cart/${cart_id}/assign" 1>/dev/null
  echo "> cart assigned to account"

  # Submit order
  order="$(request_auth POST "/order/cart/${cart_id}/checkout" '{"autoPayWithPreferredPaymentMethod":false,"waiveRetractationPeriod":false}' | $JQ_BIN -cr 'del(.contracts)')"
  $DEBUG && echo "$order" | $JQ_BIN -cr .

  order_url="$(echo "$order" | $JQ_BIN -r .url)"
  echo "> order completed url=$order_url"
}

main "$@"
