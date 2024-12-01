# Usage <img src="./assets/bash.svg" width="24">

## List available servers

```
$ kimsufi-notifier list --help
List servers from OVH Eco (including Kimsufi) catalog

Usage:
  kimsufi-notifier list [flags]

Examples:
  kimsufi-notifier list --category kimsufi
  kimsufi-notifier list --country US --endpoint ovh-us

Flags:
      --category string       category to filter on (allowed values: kimsufi, soyoustart, rise)
  -d, --datacenters strings   datacenter(s) to filter on, comma separated list (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
  -h, --help                  help for list
  -p, --plan-code string      plan code to filter on (e.g. 24ska01)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```

## Check availability

```
$ kimsufi-notifier check --help
Check OVH Eco (including Kimsufi) server availability

datacenters are the available datacenters for this plan

Usage:
  kimsufi-notifier check [flags]

Examples:
  kimsufi-notifier check --plan-code 24ska01
  kimsufi-notifier check --plan-code 24ska01 --datacenters gra,rbx

Flags:
  -d, --datacenters strings     datacenter(s) to filter on, comma separated list (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
      --help                    help for check
  -h, --human count             Human output, more h makes it better (e.g. -h, -hh)
      --list-options            list available item options
  -o, --option stringToString   options to filter on, comma separated list of key=value, see --list-options for available options (e.g. memory=ram-64g-noecc-2133) (default [])
  -p, --plan-code string        plan code name (e.g. 24ska01)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```

## Order a server

```
$ kimsufi-notifier order --help
Place an order for a servers from OVH Eco (including Kimsufi) catalog

Usage:
  kimsufi-notifier order [flags]

Examples:
  kimsufi-notifier order --plan-code 24ska01 --datacenter rbx --dry-run
  kimsufi-notifier order --plan-code 25skle01 --datacenter bhs --item-option memory=ram-32g-noecc-1333-25skle01,storage=softraid-3x2000sa-25skle01

Flags:
      --auto-pay                            automatically pay the order
  -d, --datacenter string                   datacenter (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
  -n, --dry-run                             only create a cart and do not submit the order
  -h, --help                                help for order
  -i, --item-configuration stringToString   item configuration, see --list-configurations for available values (e.g. region=europe) (default [])
  -o, --item-option stringToString          item option, see --list-options for available values (e.g. memory=ram-64g-noecc-2133-24ska01) (default [])
      --list-configurations                 list available item configurations
      --list-options                        list available item options
      --list-prices                         list available prices
      --ovh-app-key string                  environement variable name for OVH API application key (default "OVH_APP_KEY")
      --ovh-app-secret string               environement variable name for OVH API application secret (default "OVH_APP_SECRET")
      --ovh-consumer-key string             environement variable name for OVH API consumer key (default "OVH_CONSUMER_KEY")
  -p, --plan-code string                    plan code name (e.g. 24ska01)
      --price-duration string               price duration, see --list-prices for available values (default "P1M")
      --price-mode string                   price mode, see --list-prices for available values (default "default")
  -q, --quantity int                        item quantity (default 1)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```
