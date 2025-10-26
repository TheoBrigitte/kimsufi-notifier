#!/usr/bin/env bash
#
# Author: Théo Brigitte
# Date: 2025-10-26

# Usage: order.sh [options]
#
# Place an order for a servers from OVH Eco (including Kimsufi) catalog
#
# Arguments
#   -c, --country              Country code (required)
#                                Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
#                                Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS
#                                Allowed values with -e ovh-us : US
#   --datacenter               Datacenter code (default: from config when only one is set)
#                                Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)
#   -e, --endpoint             OVH API endpoint (default: $ENDPOINT)
#                                Allowed values: ovh-eu, ovh-ca, ovh-us
#   -i, --item-configuration   Item configuration in the form label=value (e.g. region=europe)
#                                use --show-configurations to list available configurations
#                                default to auto-configure when only one allowed value, otherwise ask user via prompt
#       --show-configurations  Show available configurations for the selected server
#       --item-option          Item option in the form label=value (e.g. memory=ram-64g-noecc-2133-24ska01)
#                                use --show-options to list available options
#                                default to cheapest option when multiple are available
#       --show-options         Show available options for the selected server
#   -d, --debug                Enable debug mode (default: $DEBUG)
#       --dry-run              Do not create the order, only configure the cart (default: $DRY_RUN)
#   -h, --help                 Display this help message
#   -q, --quantity             Quantity of items to order (default: $QUANTITY)
#       --price-mode           Billing price type (default: $PRICE_MODE)
#                                use --show-prices to list available price modes
#       --price-duration       Billing duration (default: $PRICE_DURATION)
#                                use --show-prices to list available price durations
#       --show-prices          Show available prices for the selected server
#
#   Arguments can also be set as environment variables see config.env.example
#   Command line arguments take precedence over environment variables
#
# Example:
#     check.sh
#     check.sh --item-configuration region=europe
#     check.sh --item-configuration region=europe --datacenter fra

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

echo_stderr() {
    >&2 echo "$@"
}

jq_stderr() {
    echo "$@" | $JQ_BIN -cr . 1>&2
}

# Helper function - prints an error message and exits
exit_error() {
    echo_stderr "Error: $1"
    exit 1
}

usage() {
  sed -Ene '/#\s?Usage:/,/^([^#]|$)/{p; /^([^#]|$)/q}' "$0" | sed -e '$d; s/#\s\?//'
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
    echo "> item configuration, select a value for $label"

    i=0
    for value in $(echo "$configuration" | $JQ_BIN -r '.allowedValues[]'); do
      echo "> $i. $value"
      i=$((i+1))
    done

    echo -n "> Choice: "
    read index
    value="$(echo "$configuration" | $JQ_BIN -r .allowedValues[$index])"
    echo "> item manual-configuration $label=$value"
    result="$(request POST "/order/cart/${cart_id}/item/${item_id}/configuration" '{"label":"'"$label"'","value":"'"$value"'"}')"
    if $DEBUG; then
      jq_stderr "$result"
    fi
  done

  exec 6<&-
}

# item_auto_option configures item with mandatory options, choosing the cheapest available
item_auto_option() {
  local cart_id="$1"
  local item_id=$2
  local plan_code="$3"
  local price_mode="$4"
  local price_duration="$5"
  shift 5
  local options_configured=("$@")

  exec 6<<<$(request GET "/order/cart/${cart_id}/eco/options?planCode=${plan_code}" | $JQ_BIN -cr '.[]')

  declare -A familyPlanCode
  declare -A familyPrices
  while read <&6 option; do
    mandatory="$(echo "$option" | $JQ_BIN -r .mandatory)"
    if [ "$mandatory" != "true" ]; then
      continue
    fi
    family="$(echo "$option" | $JQ_BIN -r .family)"
    if [[ ${options_configured[@]} =~ $family ]]; then
      continue
    fi
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

  for family in "${!familyPlanCode[@]}"; do
    option="${familyPlanCode[$family]}"
    echo_stderr "> item auto option $family=$option"
    result="$(request POST "/order/cart/${cart_id}/eco/options" '{"quantity": 1, "duration": "'"$price_duration"'", "pricingMode":"'"$price_mode"'", "planCode":"'"$option"'", "itemId": '$item_id'}')"
    if $DEBUG; then
      jq_stderr "$result"
    fi
  done
}

# item_user_option configures item with mandatory options, choosing the cheapest available
item_user_option() {
  local cart_id="$1"
  local item_id=$2
  local price_mode="$3"
  local price_duration="$4"
  shift 4
  local options=("$@")

  local keys=()
  for option in "${options[@]}"; do
    family="$(echo "$option" | cut -d= -f1)"
    planCode="$(echo "$option" | cut -d= -f2)"
    echo_stderr "> item user option $family=$planCode"
    result="$(request POST "/order/cart/${cart_id}/eco/options" '{"quantity": 1, "duration": "'"$price_duration"'", "pricingMode":"'"$price_mode"'", "planCode":"'"$planCode"'", "itemId": '$item_id'}')"
    if $DEBUG; then
      jq_stderr "$result"
    fi
    keys+=("$family")
  done

  echo "${keys[@]}"
}

print_server_configurations() {
  local cart_id="$1"
  local item_id="$2"

  exec 6<<<$(request GET "/order/cart/${cart_id}/item/${item_id}/requiredConfiguration" | $JQ_BIN -cr '.[]|select(.required==true)|select(.allowedValues|length != 1)')

  output=""
  while read <&6 configuration; do
    label="$(echo "$configuration" | $JQ_BIN -r .label)"
    exec 7<<<$(echo "$configuration" | $JQ_BIN -cr .allowedValues[])
    while read <&7 value; do
      output+="$label=$value\n"
    done
  done

  exec 6<&- 7<&-

  echo -e "$output" | column -t -N "Configuration" -o '    '
}

print_server_options() {
  cart_id="$1"
  plan_code="$2"
  price_mode="$3"
  price_duration="$4"

  exec 6<<<$(request GET "/order/cart/${cart_id}/eco/options?planCode=${plan_code}" | $JQ_BIN -cr '.[]')
  declare -A familyOptions
  declare -A familyDetails
  declare -A familyDefaults
  declare -A familyPrices
  declare -A familyLowestPrices
  while read <&6 option; do
    mandatory="$(echo "$option" | $JQ_BIN -r .mandatory)"
    if [ "$mandatory" != "true" ]; then
      continue
    fi
    family="$(echo "$option" | $JQ_BIN -r .family)"
    code="$(echo "$option" | $JQ_BIN -r .planCode)"
    name="$(echo "$option" | $JQ_BIN -r .productName)"
    price="$(echo "$option" | $JQ_BIN -r '.prices[]|select((.pricingMode == "'"$price_mode"'") and (.duration == "'"$price_duration"'"))')"
    priceText="$(echo "$price" | $JQ_BIN -r .price.text)"
    priceUcents="$(echo "$price" | $JQ_BIN -r .priceInUcents)"

    if ! [ "${familyOptions[$family]+x}" ]; then
      familyOptions[$family]="$code"
      familyDetails[$code]="$name:$mandatory:$priceText"
    else
      familyOptions[$family]="${familyOptions[$family]}:$code"
      familyDetails[$code]="$name:$mandatory:$priceText"
    fi

    if ! [ "${familyLowestPrices[$family]+x}" ]; then
      familyDefaults[$family]=$code
      familyLowestPrices[$family]=$priceUcents
    elif [ "$priceUcents" -lt "${familyLowestPrices[$family]}" ]; then
      familyDefaults[$family]=$code
      familyLowestPrices[$family]=$priceUcents
    fi
  done
  exec 6<&-

  output=""
  for key in ${!familyOptions[@]}; do
    while read option; do
      default=false
      if [ "$option" == "${familyDefaults[$key]}" ]; then
        default=true
      fi
      output+="$key=$option:${familyDetails[$option]}:$default\n"
    done <<<$(echo ${familyOptions[$key]} | tr ':' '\n')
  done
  echo -e "$output" | column -t -s ':' -N "Option,Name,Mandatory,Price,Default" -o '    '
}

print_server_prices() {
  cart_id="$1"
  plan_code="$2"

  request GET "/order/cart/${cart_id}/eco?planCode=${plan_code}" | \
    $JQ_BIN -r '.[] | select(.planCode == "'"$plan_code"'") | .prices[] | [ .duration, .pricingMode, .price.text, .description ] | @tsv' | \
    sort -k1,1 -k2n,2 -b -t $'\t' | \
    column -s $'\t' -t -N "Duration,Mode,Price,Description" -o '    '
}

main() {
  # Load configuration and common tools
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
  local item_options=()

  ARGS=$(getopt -o 'c:d:e:hi:p:q:' --long 'country:,datacenter:,item-configuration:,item-option:,debug,dry-run,endpoint:,help,quantity:,plan-code:,price-duration:,price-mode:,show-configurations,show-options,show-prices' -- "$@")
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
      --item-option)
        echo "$2" | grep -q '=' || \
          exit_error "Error: invalid item option '$2'"
        item_options+=("$2")
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
      --show-configurations)
        SHOW_CONFIGURATIONS=true
        shift 1
        continue
        ;;
      --show-options)
        SHOW_OPTIONS=true
        shift 1
        continue
        ;;
      --show-prices)
        SHOW_PRICES=true
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

  if [ -z "${COUNTRY-}" ]; then
    echo_stderr "Error: COUNTRY is not set"
    echo_stderr
    usage
    exit 3
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
  if $DEBUG; then
    jq_stderr "$cart"
  fi

  cart_id="$(echo "$cart" | $JQ_BIN -r .cartId)"
  if [ -z "$cart_id" ]; then
    echo_stderr "cart_id is empty"
    exit 3
  fi
  echo "> cart created id=$cart_id"

  if $SHOW_OPTIONS; then
    print_server_options "$cart_id" "$PLAN_CODE" "$PRICE_MODE" "$PRICE_DURATION"
    exit 0
  fi

  if $SHOW_PRICES; then
    print_server_prices "$cart_id" "$PLAN_CODE"
    exit 0
  fi

  # Add item to cart
  cart_updated="$(request POST "/order/cart/${cart_id}/eco" '{"planCode":"'"${PLAN_CODE}"'","quantity": '${QUANTITY}', "pricingMode":"'"${PRICE_MODE}"'","duration":"'${PRICE_DURATION}'"}')"
  if $DEBUG; then
    jq_stderr "$cart_updated"
  fi

  item_id="$(echo "$cart_updated" | $JQ_BIN -r .itemId)"
  if [ -z "$item_id" ]; then
    echo_stderr "> item_id is empty"
    exit 3
  fi
  echo "> cart updated with item id=$item_id"

  if $SHOW_CONFIGURATIONS; then
    print_server_configurations "$cart_id" "$item_id"
    exit 0
  fi

  # Configure item
  labels_auto_configured="$(item_auto_configuration "$cart_id" "$item_id")"
  labels_user_configured="$(item_user_configuration "$cart_id" "$item_id" "${item_configurations[@]}")"
  labels_configured=( "${labels_auto_configured[@]}" "${labels_user_configured[@]}" )
  item_manual_configuration "$cart_id" "$item_id" "${labels_configured[@]}"

  # Configure eco options
  options_configured="$(item_user_option "$cart_id" "$item_id" "$PRICE_MODE" "$PRICE_DURATION" "${item_options[@]}")"
  item_auto_option "$cart_id" $item_id "$PLAN_CODE" "$PRICE_MODE" "$PRICE_DURATION" "${options_configured[@]}"

  if $DRY_RUN; then
    echo_stderr "> dry-run enabled, skipping order completion"
    exit 0
  fi

  # Assign cart to account
  request_auth POST "/order/cart/${cart_id}/assign" 1>/dev/null
  echo "> cart assigned to account"

  # Submit order
  order="$(request_auth POST "/order/cart/${cart_id}/checkout" '{"autoPayWithPreferredPaymentMethod":false,"waiveRetractationPeriod":false}' | $JQ_BIN -cr 'del(.contracts)')"
  if $DEBUG; then
    jq_stderr "$order"
  fi

  order_url="$(echo "$order" | $JQ_BIN -r .url)"
  echo "> order completed url=$order_url"
}

main "$@"
