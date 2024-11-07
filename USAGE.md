# Usage <img src="./assets/bash.svg" width="24">

## List available servers

```
$ bin/list.sh -h
Usage: list.sh

List servers from OVH Eco (including Kimsufi) catalog

Arguments
  --category       Server category (default all)
                     Allowed values: kimsufi, soyoustart, rise, uncategorized
  -c, --country    Country code (required)
                     Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                     Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS
                     Allowed values with -e ovh-us : US
  -e, --endpoint   OVH API endpoint (default: ovh-eu)
                     Allowed values: ovh-eu, ovh-ca, ovh-us
  -d, --debug      Enable debug mode (default: false)
  -h, --help       Display this help message

  Arguments can also be set as environment variables see config.env.example
  Command line arguments take precedence over environment variables

Example:
    list.sh --country FR
    list.sh --country FR --category kimsufi
```

## Check availability

```
$ bin/check.sh -h
Usage: check.sh

Check OVH Eco (including Kimsufi) server availability

Arguments
  -p, --plan-code  Plan code to check (e.g. 24ska01)
  --datacenters    Comma-separated list of datacenters to check availability for (default all)
                     Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)
  -e, --endpoint   OVH API endpoint (default: ovh-eu)
                     Allowed values: ovh-eu, ovh-ca, ovh-us
  -d, --debug      Enable debug mode (default: false)
  -h, --help       Display this help message

  Arguments can also be set as environment variables see config.env.example
  Command line arguments take precedence over environment variables

Environment variables
    DISCORD_WEBHOOK       Webhook URL to use for Discord notification service
    GOTIFY_URL            URL to use for Gotify notification service
    GOTIFY_TOKEN          token to use for Gotify notification service
    GOTIFY_PRIORITY       prority for Gotify notification service
    OPSGENIE_API_KEY      API key for OpsGenie to receive notifications
    TELEGRAM_BOT_TOKEN    Bot token for Telegram to receive notifications
    TELEGRAM_CHAT_ID      Chat ID for Telegram to receive notifications
    HEALTHCHECKS_IO_UUID  UUID for healthchecks.io to ping after successful run

Example:
  check.sh --plan-code 24ska01
  check.sh --plan-code 24ska01 --datacenters fr,gra,rbx,sbg
```

## Order a server

```
$ bin/order.sh -h
Usage: order.sh

Place an order for a servers from OVH Eco (including Kimsufi) catalog

Arguments
  -c, --country     Country code (required)
                      Allowed values with -e ovh-eu : CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                      Allowed values with -e ovh-ca : ASIA, AU, CA, IN, QC, SG, WE, WS
                      Allowed values with -e ovh-us : US
  --datacenter      Datacenter code (default: from config when only one is set)
                      Example values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw (non exhaustive list)
  -e, --endpoint    OVH API endpoint (default: ovh-eu)
                      Allowed values: ovh-eu, ovh-ca, ovh-us
  -i, --item-configuration
                      Item configuration in the form 'label=value'
  -d, --debug       Enable debug mode (default: false)
  -h, --help        Display this help message
  -q, --quantity    Quantity of items to order (default: 1)
  --price-mode      Billing price type (default: default)
  --price-duration  Billing duration (default: P1M)

  Arguments can also be set as environment variables see config.env.example
  Command line arguments take precedence over environment variables

Example:
    order.sh
    order.sh --item-configuration region=europe
    order.sh --item-configuration region=europe --datacenter fra
```
